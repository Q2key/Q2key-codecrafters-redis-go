package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	handlers2 "github.com/codecrafters-io/redis-starter-go/app/handlers"

	"github.com/codecrafters-io/redis-starter-go/app/core"
)

func RunInstance(ctx context.Context, ins core.Redis) {
	port := ins.Config.Port

	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		log.Fatalf("\r\nFailed to bind to port %s", port)
	}

	handlers := map[string]handlers2.Handler{
		"CONFIG":   handlers2.NewConfigHandler(ins),
		"GET":      handlers2.NewGetHandler(ins),
		"SET":      handlers2.NewSetHandler(ins),
		"PING":     handlers2.NewPingHandler(ins),
		"ECHO":     handlers2.NewEchoHandler(ins),
		"KEYS":     handlers2.NewKeysHandler(ins),
		"INFO":     handlers2.NewInfoHandler(ins),
		"REPLCONF": handlers2.NewReplConfHandler(ins),
		"PSYNC":    handlers2.NewPsyncHandler(ins),
		"WAIT":     handlers2.NewWaitHandler(ins),
		"TYPE":     handlers2.NewTypeHandler(ins),
	}

	go ins.InitHandshakeWithMaster()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("Something went wrong with tcp connection...")
		}

		go handleRedisConnection(conn, ins, handlers)
	}
}

func handleRedisConnection(conn net.Conn, ins core.Redis, handlers map[string]handlers2.Handler) {
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
		if err != nil || len(args) == 0 {
			// log.Fatal("error reading from connection")
		}

		name := args[0]
		isWrite := name == "SET"

		h := handlers[name]
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
