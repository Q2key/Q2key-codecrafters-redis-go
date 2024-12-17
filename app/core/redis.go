package core

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/codecrafters-io/redis-starter-go/app/contracts"
	"github.com/codecrafters-io/redis-starter-go/app/handlers"
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

type Redis struct {
	Config      Config
	Store       Store
	RepConnPool *map[string]contracts.Connector
	MasterConn  contracts.Connector
	AckChan     *chan Ack
	Bytes       *int
	Handlers    map[string]contracts.Handler
}

func NewRedis(_ context.Context, config Config) *Redis {
	ch := make(chan Ack)

	bytes := 0
	ins := &Redis{
		Store:       *NewStore(),
		Config:      config,
		RepConnPool: &map[string]contracts.Connector{},
		AckChan:     &ch,
		Bytes:       &bytes,
	}

	ins.Handlers = map[string]contracts.Handler{
		"CONFIG":   handlers.NewConfigHandler(*ins),
		"GET":      handlers.NewGetHandler(*ins),
		"SET":      handlers.NewSetHandler(*ins),
		"PING":     handlers.NewPingHandler(*ins),
		"ECHO":     handlers.NewEchoHandler(*ins),
		"KEYS":     handlers.NewKeysHandler(*ins),
		"INFO":     handlers.NewInfoHandler(*ins),
		"REPLCONF": handlers.NewReplConfHandler(*ins),
		"PSYNC":    handlers.NewPsyncHandler(*ins),
		"WAIT":     handlers.NewWaitHandler(*ins),
		"TYPE":     handlers.NewTypeHandler(*ins),
	}

	ins.TryReadDb()

	return ins
}

func (r *Redis) HandleRedisConnection(conn net.Conn) {
	defer conn.Close()
	buff := make([]byte, 215)

	redisCon := NewRConn(&conn)
	for {
		n, err := conn.Read(buff)
		if err == io.EOF {
			continue
		}

		payload := buff[:n]

		_, args := FromRedisStringToStringArray(string(payload))

		name := args[0]
		isWrite := name == "SET"

		h := r.Handlers[name]
		h.Handle(redisCon, args, &payload)

		if r.Config.IsMaster() && isWrite {
			r.SendToReplicas(&payload)
		}
	}
}

func (r *Redis) InitHandshakeWithMaster() {
	rep := r.Config.Replica
	if rep == nil {
		return
	}

	host, port := rep.OriginHost, rep.OriginPort
	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		log.Fatal("Handshake connection error")
	}

	masterConn := NewRConn(&conn)

	r.RegisterMasterConn(masterConn)
	defer conn.Close()

	if conn == nil {
		return
	}

	// Handshake 1
	conn.Write([]byte("*1\r\n$4\r\nPING\r\n"))

	buff := make([]byte, 512)
	n, _ := conn.Read(buff)
	if string(buff[:n]) != "+PONG\r\n" {
		return
	}

	reader := bufio.NewReader(conn)

	// Handshake 2.1
	req := StringsToRedisStrings([]string{"REPLCONF", "listening-port", r.Config.Port})
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
			continue
		}

		sbuf := buf[:n]

		res.Write(sbuf)
		if bytes.Contains(sbuf, repb) {
			lx := res.Len() - shift
			req = StringsToRedisStrings([]string{"REPLCONF", "ACK", strconv.Itoa(lx - rshift)})
			conn.Write([]byte(req))
		}

		payload := r.bytesToCommandMap(res.Bytes())
		for k, v := range payload {
			r.Store.Set(k, v.Value)
		}
	}
}

func (r *Redis) TryReadDb() {
	if r.Config.DbFileName == "" || r.Config.Dir == "" {
		return
	}

	path := fmt.Sprintf("%s/%s", r.Config.Dir, r.Config.DbFileName)

	db := NewRedisDB(path)
	if !db.IsFileExists(r.Config.DbFileName) {
		_ = os.Mkdir(r.Config.Dir, os.ModeDir)
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
		r.Store.Set(k, v)
		exp, ok := db.GetExpires()[k]
		if ok {
			r.Store.SetExpiredAt(k, exp)
		}
	}
}

func (r *Redis) SendToReplicas(buff *[]byte) {
	*r.Bytes += len(*buff)
	for _, r := range *r.RepConnPool {
		r.Conn().Write(*buff)
	}
}

func (r *Redis) RegisterReplicaConn(conn contracts.Connector) {
	(*r.RepConnPool)[(conn).Id()] = conn
}

func (r *Redis) RegisterMasterConn(conn contracts.Connector) {
	r.MasterConn = conn
}

func (r *Redis) UpdateReplica(id string, offset int) {
	rep := (*r.RepConnPool)[id]
	rep.SetOffset(offset)
	(*r.RepConnPool)[id] = rep
}

func (r *Redis) GetWrittenBytes() int {
	return *r.Bytes
}

func (r *Redis) bytesToCommandMap(buf []byte) map[string]StoreValue {
	res := map[string]StoreValue{}

	j := 0
	for i, ch := range buf {
		if string(ch) == "*" {
			j = i
			break
		}
	}

	_, arr := FromRedisStringToStringArray(string(buf)[j:])
	for i, v := range arr {
		if v == "SET" && i+2 <= len(arr) {
			res[arr[i+1]] = StoreValue{
				Value: arr[i+2],
			}
		}
	}

	return res
}
