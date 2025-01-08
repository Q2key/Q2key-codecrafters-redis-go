package core

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

type Master struct {
	Redis
	RepConnPool   *map[string]Conn
	AckChan       *chan Ack
	ReceivedBytes *int
}

func NewMaster(_ context.Context, config Config) *Master {
	ch := make(chan Ack)

	receivedBytes := 0
	ins := &Master{
		Redis: Redis{
			Store:  NewStore(),
			Config: config,
			Commands: map[string]CommandHandler{
				"PING":     handlePing,
				"INFO":     handleInfo,
				"KEYS":     handleKeys,
				"GET":      handleGet,
				"SET":      handleSet,
				"CONFIG":   handleConfig,
				"XADD":     handleXadd,
				"XRANGE":   handleXrange,
				"ECHO":     handleEcho,
				"WAIT":     handleWaitAsMaster,
				"REPLCONF": handleReplconf,
				"PSYNC":    handlePsync,
				"TYPE":     handleType,
			},
		},
		RepConnPool:   &map[string]Conn{},
		AckChan:       &ch,
		ReceivedBytes: &receivedBytes,
	}

	return ins
}

func (r *Master) Init() {
	r.TryReadDb()
}

func (r *Master) GetConfig() *Config {
	return &r.Config
}

func (r *Master) GetStore() *Store {
	return &r.Store
}

func (r *Master) HandleTCP(conn net.Conn) {
	defer conn.Close()
	buff := make([]byte, 215)

	redisCon := NewRConn(&conn)

	for {
		n, err := conn.Read(buff)
		if err == io.EOF {
			continue
		}

		payload := buff[:n]

		args := FromRedisStringToStringArray(string(payload))
		if len(args) == 0 {
			return
		}

		command := args[0]
		handler, ok := r.Commands[command]
		if !ok {
			continue
		}

		handler(r, *redisCon, args)

		if IsWriteCommand(command) {
			r.SendToReplicas(&payload)
		}
	}
}

func (r *Master) SendToReplicas(buff *[]byte) {
	*r.ReceivedBytes += len(*buff)
	for _, r := range *r.RepConnPool {
		_, err := r.Conn().Write(*buff)
		if err != nil {
			return
		}
	}
}

func (r *Master) RegisterReplicaConn(conn Conn) {
	(*r.RepConnPool)[(conn).Id()] = conn
}

func (r *Master) UpdateReplica(id string, offset int) {
	rep := (*r.RepConnPool)[id]
	rep.SetOffset(offset)
	(*r.RepConnPool)[id] = rep
}

func (r *Master) TryReadDb() {
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
		r.Store.Set(k, v, STRING)
		exp, ok := db.GetExpires()[k]
		if ok {
			r.Store.SetExpiredAt(k, exp)
		}
	}
}
