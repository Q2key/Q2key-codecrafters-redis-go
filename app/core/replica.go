package core

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"log"
	"net"
	"strconv"
)

type Replica struct {
	Redis
	MasterConn    Conn
	AckChan       *chan Ack
	ReceivedBytes *int
}

func NewReplica(_ context.Context, config Config) *Replica {
	ch := make(chan Ack)
	receivedBytes := 0
	ins := &Replica{
		Redis: Redis{
			Store:  NewStore(),
			Config: config,
			Commands: map[string]CommandHandler{
				"PING":   handlePing,
				"INFO":   handleInfo,
				"KEYS":   handleKeys,
				"GET":    handleGet,
				"CONFIG": handleConfig,
				"ECHO":   handleEcho,
				"TYPE":   handleType,
			},
		},
		AckChan:       &ch,
		ReceivedBytes: &receivedBytes,
	}

	return ins
}

func (r *Replica) Init() {
	go r.Handshake()
}

func (r *Replica) GetConfig() *Config {
	return &r.Config
}

func (r *Replica) GetStore() *Store {
	return &r.Store
}

func (r *Replica) HandleTCP(conn net.Conn) {
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
	}
}

func (r *Replica) RegisterMasterConn(conn Conn) {
	r.MasterConn = conn
}

func (r *Replica) Handshake() {
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

	r.RegisterMasterConn(*masterConn)
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
	req := ToRedisStrings([]string{"REPLCONF", "listening-port", r.Config.Port})
	conn.Write([]byte(req))
	reader.ReadBytes('\n')

	// Handshake 2.2
	req = ToRedisStrings([]string{"REPLCONF", "capa", "psync2"})
	conn.Write([]byte(req))
	reader.ReadBytes('\n')

	// Handshake 3
	req = ToRedisStrings([]string{"PSYNC", "?", "-1"})

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
			req = ToRedisStrings([]string{"REPLCONF", "ACK", strconv.Itoa(lx - rshift)})
			conn.Write([]byte(req))
		}

		payload := r.Store.BytesToCommandMap(res.Bytes())
		for k, v := range payload {
			r.Store.Set(k, v.Value, STRING)
		}
	}
}
