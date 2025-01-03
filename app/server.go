package main

import (
	"context"
	"fmt"
	"github.com/codecrafters-io/redis-starter-go/app/core"
	"github.com/codecrafters-io/redis-starter-go/app/handlers"
	"io"
	"log"
	"net"
	"os"
)

func RunInstance(ctx context.Context, ins core.Redis) {
	port := ins.Config.Port

	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		log.Fatalf("\r\nFailed to bind to port %s", port)
	}

	hs := map[string]handlers.Handler{
		"CONFIG":   handlers.NewConfigHandler(ins),
		"GET":      handlers.NewGetHandler(ins),
		"SET":      handlers.NewSetHandler(ins),
		"PING":     handlers.NewPingHandler(ins),
		"ECHO":     handlers.NewEchoHandler(ins),
		"KEYS":     handlers.NewKeysHandler(ins),
		"INFO":     handlers.NewInfoHandler(ins),
		"REPLCONF": handlers.NewReplConfHandler(ins),
		"PSYNC":    handlers.NewPsyncHandler(ins),
		"WAIT":     handlers.NewWaitHandler(ins),
		"TYPE":     handlers.NewTypeHandler(ins),
		"XADD":     handlers.NewXaddHandler(ins),
	}

	go ins.InitHandshakeWithMaster()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("Something went wrong with tcp connection...")
		}

		go handleRedisConnection(conn, ins, hs)
	}
}

func handleRedisConnection(conn net.Conn, ins core.Redis, hs map[string]handlers.Handler) {
	defer conn.Close()
	buff := make([]byte, 215)

	redisCon := core.NewRConn(&conn)
	for {
		n, err := conn.Read(buff)
		if err == io.EOF {
			continue
		}

		payload := buff[:n]

		err, args := core.FromRedisStringToStringArray(string(payload))

		name := args[0]
		isWrite := name == "SET"

		h := hs[name]

		h.Handle(*redisCon, args, &payload)

		if ins.Config.IsMaster() && isWrite {
			ins.SendToReplicas(&payload)
		}
	}
}

func main() {
	ctx := context.Background()
	config := core.NewConfig().FromArguments(os.Args)
	redis := core.NewRedis(ctx, *config)
	RunInstance(ctx, *redis)
}
