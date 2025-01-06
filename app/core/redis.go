package core

import (
	"context"
	"io"
	"log"
	"net"
)

type Redis struct {
	Config        Config
	Store         Store
	RepConnPool   *map[string]Conn
	MasterConn    Conn
	AckChan       *chan Ack
	ReceivedBytes *int
}

func NewRedis(_ context.Context, config Config) *Redis {
	ch := make(chan Ack)

	receivedBytes := 0
	ins := &Redis{
		Store:         *NewStore(),
		Config:        config,
		RepConnPool:   &map[string]Conn{},
		AckChan:       &ch,
		ReceivedBytes: &receivedBytes,
	}

	ins.TryReadDb()

	return ins
}

func (r *Redis) RegisterMasterConn(conn Conn) {
	r.MasterConn = conn
}

func (r *Redis) HandleTCP(conn net.Conn) {
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
		switch command {
		case "PING":
			handlePING(*redisCon)
		case "INFO":
			handleINFO(r, *redisCon, args)
		case "KEYS":
			handleKEYS(r, *redisCon, args)
		case "GET":
			handleGET(r, *redisCon, args)
		case "SET":
			handleSET(r, *redisCon, args)
		case "ECHO":
			handleECHO(*redisCon, args)
		case "REPLCONF":
			handleREPLCONF(r, *redisCon, args)
		case "CONFIG":
			handleCONFIG(r, *redisCon, args)
		case "PSYNC":
			handlePSYNC(r, *redisCon)
		case "WAIT":
			handleWAIT(r, *redisCon, args)
		case "TYPE":
			handleTYPE(r, *redisCon, args)
		case "XADD":
			handleXADD(r, *redisCon, args)
		default:
			log.Fatal("Unknown command")
		}

		if r.Config.IsMaster() && IsWriteCommand(command) {
			r.SendToReplicas(&payload)
		}
	}
}
