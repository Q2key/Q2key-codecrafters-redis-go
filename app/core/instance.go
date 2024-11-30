package core

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/client"
	"github.com/codecrafters-io/redis-starter-go/app/contracts"
)

type Instance struct {
	ReplicaId   string
	Config      contracts.Config
	Store       contracts.Store
	RepConnPool map[string]*net.Conn
	MasterConn  *net.Conn
	Ch          chan bool
	mu          sync.Mutex
}

const FakeReplicaId = "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb"

func NewRedisInstance(config contracts.Config) *Instance {
	ch := make(chan bool, 1)
	ins := &Instance{
		ReplicaId:   FakeReplicaId,
		Store:       contracts.Store{},
		Config:      config,
		RepConnPool: map[string]*net.Conn{},
		Ch:          ch,
		mu:          sync.Mutex{},
	}

	if !config.IsMaster() {
		go ins.HandShakeMaster()
	}

	ins.TryReadDb()

	return ins
}

func (r *Instance) HandShakeMaster() {
	if r.GetConfig().IsMaster() {
		return
	}

	fmt.Println("Starting handshake")

	rep := r.Config.GetReplica()
	if rep == nil {
		return
	}

	host, port := rep.OriginHost, rep.OriginPort
	tcp := client.NewTcpClient(host, port)

	err := tcp.Connect()
	if err != nil {
		fmt.Println(err.Error())
	}

	r.RegisterMasterConn(*tcp.Conn())

	// Handshake 1
	buff := send("*1\r\n$4\r\nPING\r\n", tcp)
	if string(*buff) != "+PONG\r\n" {
		return
	}

	// Handshake 2
	req := FromStringArrayToRedisStringArray([]string{"REPLCONF", "listening-port", r.Config.GetPort()})
	send(req, tcp)

	req = FromStringArrayToRedisStringArray([]string{"REPLCONF", "capa", "psync2"})
	send(req, tcp)
	// Handshake 3
	req = FromStringArrayToRedisStringArray([]string{"PSYNC", "?", "-1"})
	send(req, tcp)
}

func send(req string, client *client.TcpClient) *[]byte {
	buff, err := (*client).SendBytes([]byte(req))
	if err != nil {
		fmt.Println(err)
	}

	return buff
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
	// r.mu.Lock()
	v := r.Store[key]
	// r.mu.Unlock()
	return v
}

func (r *Instance) GetReplicaId() string {
	return r.ReplicaId
}

func (r *Instance) Set(key string, value string) {
	// r.mu.Lock()
	r.Store[key] = &InstanceValue{
		Value: value,
	}
	// r.mu.Unlock()
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

func (r *Instance) Propagate(buff []byte) {
	if r.MasterConn != nil {
		_, err := (*r.MasterConn).Write(buff)
		if err != nil {
			fmt.Print(err)
		}
	}
}

func (r *Instance) GetMasterConn() *net.Conn {
	return r.MasterConn
}

func (r *Instance) GetChan() chan bool {
	return r.Ch
}
