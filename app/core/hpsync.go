package core

import (
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	MockDbContent = "524544495330303131fa0972656469732d76657205372e322e30fa0a72656469732d62697473c040fa056374696d65c26d08bc65fa08757365642d6d656dc2b0c41000fa08616f662d62617365c000fff06e3bfec0ff5aa2"
)

func handlePsync(r RedisInstance, conn Conn, _ []string) {
	master, ok := r.(*Master)
	if !ok {
		return
	}

	h := master
	h.RegisterReplicaConn(conn)

	resp := ToRedisSimpleString(fmt.Sprintf("FULLRESYNC %s 0", conn.Id()))
	rdbBuff, err := hex.DecodeString(MockDbContent)
	if err != nil {
		log.Fatal(err)
	}

	var sb strings.Builder
	defer sb.Reset()
	sb.WriteString(resp)
	sb.WriteString("$88\r\n")
	sb.Write(rdbBuff)

	RespondString(conn, sb.String())
}

func handleWaitAsMaster(r RedisInstance, conn Conn, args []string) {
	master, ok := r.(*Master)
	if !ok {
		return
	}

	h := master
	rep, err := strconv.Atoi(args[1])
	if err != nil {
		return
	}

	timeout, err := strconv.Atoi(args[2])
	if err != nil {
		return
	}

	done := map[string]bool{}

	for _, r := range *h.RepConnPool {
		if r.Offset() >= *h.ReceivedBytes {
			done[r.Id()] = true
			continue
		}

		go RespondString(r, "*3\r\n$8\r\nREPLCONF\r\n$6\r\nGETACK\r\n$1\r\n*\r\n")

	}

	ch := *h.AckChan

awaitingLoop:
	for len(done) < rep {
		select {
		case c := <-ch:
			{
				if c.Offset >= *h.ReceivedBytes {
					done[c.ConnId] = true
				}
			}
		case <-time.After(time.Duration(timeout) * time.Millisecond):
			{
				break awaitingLoop
			}
		}
	}

	RespondString(conn, ToRedisInteger(strconv.Itoa(len(done))))
}
