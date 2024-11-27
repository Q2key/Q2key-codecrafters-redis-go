package core

import (
	"fmt"
	"github.com/codecrafters-io/redis-starter-go/app/client"
	"github.com/codecrafters-io/redis-starter-go/app/contracts"
	"log"
	"net"
	"os"
	"time"
)

type Instance struct {
	ReplicaId     string
	Config        contracts.Config
	store         contracts.Store
	remoteAddress string
	conn          net.Conn
	repConnPool   map[string]*net.Conn
}

const FakeReplicaId = "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb"

func NewRedisInstance(config contracts.Config) *Instance {
	ins := &Instance{
		ReplicaId:   FakeReplicaId,
		store:       contracts.Store{},
		Config:      config,
		repConnPool: map[string]*net.Conn{},
	}

	ins.TryReadDb()
	if !config.IsMaster() {
		ins.HandShakeMaster()
	}

	return ins
}

func (r *Instance) HandShakeMaster() {
	rep := r.Config.GetReplica()
	if rep == nil {
		return
	}

	fmt.Println("Replica is on: " + rep.OriginPort)
	host, port := rep.OriginHost, rep.OriginPort
	tcp := client.NewTcpClient(host, port)

	err := tcp.Connect()
	if err != nil {
		log.Fatal(err)
	}

	//Handshake 1
	bytes, err := tcp.SendBytes([]byte("*1\r\n$4\r\nPING\r\n"))

	if string(*bytes) != "+PONG\r\n" {
		return
	}

	//Handshake 2
	req := FromStringArrayToRedisStringArray([]string{"REPLCONF", "listening-port", r.Config.GetPort()})

	sentBytes, err := tcp.SendBytes([]byte(req))
	if err != nil || len(*sentBytes) == 0 {
		//return
	}

	req = FromStringArrayToRedisStringArray([]string{"REPLCONF", "capa", "psync2"})
	sentBytes, err = tcp.SendBytes([]byte(req))
	if err != nil || len(*sentBytes) == 0 {
		//return
	}
	//Handshake 3
	req = FromStringArrayToRedisStringArray([]string{"PSYNC", "?", "-1"})
	sentBytes, err = tcp.SendBytes([]byte(req))
	if err != nil || len(*sentBytes) == 0 {
		//return
	}

	//tcp.Disconnect()
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
	return r.store[key]
}

func (r *Instance) GetReplicaId() string {
	return r.ReplicaId
}

func (r *Instance) Set(key string, value string) {
	r.store[key] = &InstanceValue{
		Value: value,
	}
}

func (r *Instance) GetKeys(token string) []string {
	res := make([]string, 0)
	switch token {
	case "*":
		for k := range r.store {
			res = append(res, k)
		}
	}
	return res
}

func (r *Instance) GetStore() *map[string]contracts.Value {
	return &r.store
}

func (r *Instance) SetExpiredAt(key string, expired uint64) {
	tm := GetDateFromTimeStamp(expired)
	val, ok := r.store[key]
	if ok {
		val.SetExpired(tm)
	}

	r.store[key] = val
}

func (r *Instance) SetExpiredIn(key string, expiredIn uint64) {
	exp := time.Now().UTC().Add(time.Duration(expiredIn) * time.Millisecond)
	val, ok := r.store[key]
	if ok {
		val.SetExpired(exp)
	}
	r.store[key] = val
}

func (r *Instance) GetConfig() contracts.Config {
	return r.Config
}

func (r *Instance) Replicate(buff []byte) {
	for _, c := range r.repConnPool {
		if c != nil {
			(*c).Write(buff)
		}
	}
}

func (r *Instance) RegisterReplicaConn(conn net.Conn) {
	r.repConnPool[conn.RemoteAddr().String()] = &conn
}
