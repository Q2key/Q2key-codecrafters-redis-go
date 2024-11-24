package main

import (
	"fmt"
	"github.com/codecrafters-io/redis-starter-go/app/commands"
	"github.com/codecrafters-io/redis-starter-go/app/contracts"
	"github.com/codecrafters-io/redis-starter-go/app/core"
	"github.com/codecrafters-io/redis-starter-go/app/handlers"
	"log"
	"net"
	"os"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var (
	_ = net.Listen
	_ = os.Exit
)

func RunInstance(ins contracts.Instance) {
	port := ins.GetConfig().GetPort()

	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		log.Fatalf("\r\nFailed to bind to port %s", port)
	}

	configHandler := handlers.NewConfigHandler(ins)
	getHandler := handlers.NewGetHandler(ins)
	setHandler := handlers.NewSetHandler(ins)
	pingHandler := handlers.NewPingHandler(ins)
	echoHandler := handlers.NewEchoHandler(ins)
	keysHandler := handlers.NewKeysHandler(ins)
	infoHandler := handlers.NewInfoHandler(ins)
	replconfHandler := handlers.NewReplConfHandler(ins)
	psyncHandler := handlers.NewPsyncHandler(ins)

	for {
		conn, _ := ln.Accept()
		//todo: implement handle connection
		go func() {
			buff := make([]byte, 1024*8)
			defer conn.Close()
			for {
				read, err := conn.Read(buff)
				if read == 0 {
					continue
				}

				if err != nil {
					return
				}

				err, cmd := commands.ParseCommand(string(buff))
				if err != nil {
					handlers.HandleError(conn, err)
					continue
				}

				switch cmd.Name() {
				case "CONFIG":
					configHandler.Handle(conn, cmd)
				case "GET":
					getHandler.Handle(conn, cmd)
				case "SET":
					setHandler.Handle(conn, cmd)
				case "ECHO":
					echoHandler.Handle(conn, cmd)
				case "PING":
					pingHandler.Handle(conn, cmd)
				case "KEYS":
					keysHandler.Handle(conn, cmd)
				case "INFO":
					infoHandler.Handle(conn, cmd)
				case "REPLCONF":
					replconfHandler.Handle(conn, cmd)
				case "PSYNC":
					psyncHandler.Handle(conn, cmd)
				}
			}
		}()
	}
}

func main() {
	fmt.Println("Starting server...")
	config := core.NewConfig().FromArguments(os.Args)
	redis := core.NewRedisInstance(*config)
	RunInstance(redis)
}
