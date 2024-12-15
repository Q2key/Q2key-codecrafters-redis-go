package core

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/contracts"
)

type Instance struct {
	Config      contracts.Config
	Store       contracts.Store
	RepConnPool *map[string]contracts.RedisConn
	MasterConn  contracts.RedisConn
	Scheduler   *contracts.Scheduler
	Chan        *chan contracts.Ack
	bytes       int
}

func NewInstance(_ context.Context, config contracts.Config) *Instance {
	ch := make(chan contracts.Ack)
	ins := &Instance{
		Store:       NewStore(),
		Config:      config,
		RepConnPool: &map[string]contracts.RedisConn{},
		Chan:        &ch,
	}

	if config.IsMaster() {
		waitMs := 0
		ins.Scheduler = &contracts.Scheduler{
			WaitTimeoutMS:      &waitMs,
			TotalReplicasCount: 0,
		}
	}

	ins.TryReadDb()

	return ins
}

func (r *Instance) InitHandshakeWithMaster() {
	rep := r.Config.GetReplica()
	if rep == nil {
		return
	}

	host, port := rep.OriginHost, rep.OriginPort
	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		log.Fatal("Handshake connection error")
	}

	masterConn := NewReplicMasterConn(&conn)
	r.RegisterMasterConn(&masterConn)
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
	req := StringsToRedisStrings([]string{"REPLCONF", "listening-port", r.Config.GetPort()})
	conn.Write([]byte(req))
	reader.ReadBytes('\n')

	// Handshake 2.2
	req = StringsToRedisStrings([]string{"REPLCONF", "capa", "psync2"})
	conn.Write([]byte(req))
	reader.ReadBytes('\n')

	// Handshake 3
	req = StringsToRedisStrings([]string{"PSYNC", "?", "-1"})

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
			req = StringsToRedisStrings([]string{"REPLCONF", "ACK", strconv.Itoa(lx - rshift)})
			conn.Write([]byte(req))
		}

		payload := BytesToCommandMap(res.Bytes())
		for k, v := range payload {
			r.GetStore().Set(k, v.Value)
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
		r.GetStore().Set(k, v)
		exp, ok := db.GetExpires()[k]
		if ok {
			r.GetStore().SetExpiredAt(k, exp)
		}
	}
}

func (r *Instance) GetAckChan() *chan contracts.Ack {
	return r.Chan
}

func (r *Instance) GetStore() contracts.Store {
	return r.Store
}

func (r *Instance) GetConfig() contracts.Config {
	return r.Config
}

func (r *Instance) SendToReplicas(buff *[]byte) {
	r.bytes += len(*buff)
	for _, r := range r.GetReplicas() {
		r.GetConn().Write((*buff))
	}
}

func (r *Instance) RegisterReplicaConn(conn *contracts.RedisConn) {
	(*r.RepConnPool)[(*conn).GetId()] = *conn
	r.Scheduler.IncreasTotalReplicasCounter()
}

func (r *Instance) RegisterMasterConn(conn *contracts.RedisConn) {
	r.MasterConn = *conn
}

func (r *Instance) GetMasterConn() contracts.RedisConn {
	return r.MasterConn
}

func (r *Instance) GetReplicas() map[string]contracts.RedisConn {
	return *r.RepConnPool
}

func (r *Instance) UpdateReplica(id string, offset int) {
	(*r.RepConnPool)[id].SetOffset(offset)
}

func (r *Instance) GetWrittenBytes() int {
	return r.bytes
}
