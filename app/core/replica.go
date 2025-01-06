package core

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net"
	"strconv"
)

func (r *Redis) Handshake() {
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
