package core

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/contracts"
)

type Instance struct {
	ReplicaId   string
	Config      contracts.Config
	Store       contracts.Store
	RepConnPool map[string]*net.Conn
	MasterConn  *net.Conn
}

const FakeReplicaId = "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb"

func NewRedisInstance(config contracts.Config) *Instance {
	ins := &Instance{
		ReplicaId:   FakeReplicaId,
		Store:       contracts.Store{},
		Config:      config,
		RepConnPool: map[string]*net.Conn{},
	}

	ins.TryReadDb()

	return ins
}

func (r *Instance) HandShakeMaster(ch chan []byte) {
	if r.GetConfig().IsMaster() {
		return
	}

	rep := r.Config.GetReplica()
	if rep == nil {
		return
	}

	host, port := rep.OriginHost, rep.OriginPort
	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		log.Fatal("Handshake connection error")
	}

	r.RegisterMasterConn(conn)
	defer conn.Close()

	if conn == nil {
		return
	}

	// Handshake 1
	conn.Write([]byte("*1\r\n$4\r\nPING\r\n"))

	buff := make([]byte, 512)
	n, _ := conn.Read(buff)
	if string(buff[:n]) != "+PONG\r\n" {
		fmt.Print("Expected to get PONG from master")
		return
	}

	reader := bufio.NewReader(conn)

	// Handshake 2.1
	req := FromStringArrayToRedisStringArray([]string{"REPLCONF", "listening-port", r.Config.GetPort()})
	conn.Write([]byte(req))
	reader.ReadBytes('\n')

	// Handshake 2.2
	req = FromStringArrayToRedisStringArray([]string{"REPLCONF", "capa", "psync2"})
	conn.Write([]byte(req))
	reader.ReadBytes('\n')

	// Handshake 3
	req = FromStringArrayToRedisStringArray([]string{"PSYNC", "?", "-1"})
	conn.Write([]byte(req))
	reader.ReadBytes('\n')

	// Handshake 4
	var res bytes.Buffer
	buf := make([]byte, 512)
	shift := len([]byte("$88\r\n")) + 88

	repq := "*3\r\n$8\r\nREPLCONF\r\n$6\r\nGETACK\r\n$1\r\n*\r\n"
	repb := []byte(repq)
	rshift := len(repb)

	for {
		n, err := reader.Read(buf)

		if err == io.EOF {
			break
		}

		sbuf := buf[:n]

		res.Write(sbuf)

		if bytes.Contains(sbuf, repb) {
			lx := res.Len() - shift
			req = FromStringArrayToRedisStringArray([]string{"REPLCONF", "ACK", strconv.Itoa(lx - rshift)})
			conn.Write([]byte(req))
		}

		cmap := FromBytesArrayToSetCommandMap(res.Bytes())
		for key, ival := range cmap {
			r.Set(key, ival.Value)
		}
	}
}

func (r *Instance) TryReadDb() {
	if r.Config.GetDbFileName() == "" || r.Config.GetDir() == "" {
		return
	}

	path := fmt.Sprintf("%s/%s", r.Config.GetDir(), r.Config.GetDbFileName())

	db := NewRedisDB(path)
	if !db.IsFileExists(r.Config.GetDbFileName()) {
		_ = os.Mkdir(r.Config.GetDir(), os.ModeDir)
	}

	if !db.IsFileExists(path) {
		err := db.Create()
		if err != nil {
			log.Fatal(err)
		}
	}

	err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range db.GetData() {
		r.Set(k, v)
		exp, ok := db.GetExpires()[k]
		if ok {
			r.SetExpiredAt(k, exp)
		}
	}
}

func (r *Instance) Get(key string) contracts.Value {
	v := r.Store[key]
	return v
}

func (r *Instance) GetReplicaId() string {
	return r.ReplicaId
}

func (r *Instance) Set(key string, value string) {
	r.Store[key] = &InstanceValue{
		Value: value,
	}
}

func (r *Instance) GetKeys(token string) []string {
	res := make([]string, 0)
	switch token {
	case "*":
		for k := range r.Store {
			res = append(res, k)
		}
	}
	return res
}

func (r *Instance) GetStore() *map[string]contracts.Value {
	return &r.Store
}

func (r *Instance) SetExpiredAt(key string, expired uint64) {
	tm := GetDateFromTimeStamp(expired)
	val, ok := r.Store[key]
	if ok {
		val.SetExpired(tm)
	}

	r.Store[key] = val
}

func (r *Instance) SetExpiredIn(key string, expiredIn uint64) {
	exp := time.Now().UTC().Add(time.Duration(expiredIn) * time.Millisecond)
	val, ok := r.Store[key]
	if ok {
		val.SetExpired(exp)
	}
	r.Store[key] = val
}

func (r *Instance) GetConfig() contracts.Config {
	return r.Config
}

func (r *Instance) Replicate(buff []byte) {
	for _, c := range r.RepConnPool {
		if c != nil {
			(*c).Write(buff)
		}
	}
}

func (r *Instance) RegisterReplicaConn(conn net.Conn) {
	r.RepConnPool[fmt.Sprintf("%p", conn)] = &conn
}

func (r *Instance) RegisterMasterConn(conn net.Conn) {
	r.MasterConn = &conn
}

func (r *Instance) GetMasterConn() *net.Conn {
	return r.MasterConn
}

func (r *Instance) GetReplicas() map[string]*net.Conn {
	return r.RepConnPool
}

// todo remove or refactor
func readFromMaster(conn *net.Conn, ch chan []byte) {
	if conn == nil {
		return
	}

	defer (*conn).Close()
	out := make([]byte, 512)
	for {
		i, err := (*conn).Read(out)
		if err == io.EOF {
			break
		}

		fmt.Print("HELLO")

		ch <- out[:i]
	}
}
