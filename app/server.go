package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/commands"
	"github.com/codecrafters-io/redis-starter-go/app/contracts"
	"github.com/codecrafters-io/redis-starter-go/app/core"
	"github.com/codecrafters-io/redis-starter-go/app/handlers"
)

func RunInstance(ins contracts.Instance) {
	port := ins.GetConfig().GetPort()

	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		log.Fatalf("\r\nFailed to bind to port %s", port)
	}

	mch := make(chan []byte)
	handlers := map[string]contracts.Handler{
		"CONFIG":   handlers.NewConfigHandler(ins),
		"GET":      handlers.NewGetHandler(ins),
		"SET":      handlers.NewSetHandler(ins),
		"PING":     handlers.NewPingHandler(ins),
		"ECHO":     handlers.NewEchoHandler(ins),
		"KEYS":     handlers.NewKeysHandler(ins),
		"INFO":     handlers.NewInfoHandler(ins),
		"REPLCONF": handlers.NewReplConfHandler(ins),
		"PSYNC":    handlers.NewPsyncHandler(ins),
	}

	go ins.HandShakeMaster(mch)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("Something went wrong with tcp connection...")
		}

		go handleRedisConnection(conn, ins, handlers)
	}
}

func handleRedisConnection(conn net.Conn, ins contracts.Instance, handlers map[string]contracts.Handler) {
	defer conn.Close()

	buff := make([]byte, 512)
	for {
		n, err := conn.Read(buff)
		if err != nil {
			if err == io.EOF {
				continue
			}
		}

		err, cmd := commands.ParseCommand(string(buff[:n]))
		if err != nil {
			continue
		}

		h, ok := handlers[cmd.Name()]
		if !ok {
			break
		}

		h.Handle(conn, cmd)

		if cmd.IsWrite() && ins.GetConfig().IsMaster() {
			repData := core.FromStringArrayToRedisStringArray(cmd.Args())
			ins.Replicate([]byte(repData))
		}
	}
}

func main() {
	fmt.Println("Starting server...")
	config := core.NewConfig().FromArguments(os.Args)
	redis := core.NewRedisInstance(*config)
	RunInstance(redis)
}
