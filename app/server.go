package main

import (
	"fmt"
	"github.com/codecrafters-io/redis-starter-go/app/command"
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

func main() {
	fmt.Println("Logs from your program will appear here!")
	s := core.NewRedisInstanceWithArgs(os.Args)

	s.GetConfig().GetPort()
	port := s.GetConfig().GetPort()

	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		fmt.Printf("\r\nFailed to bind to port %s", port)
		os.Exit(1)
	}

	//defer ln.Close()

	configHandler := handlers.NewConfigHandler(s)
	getHandler := handlers.NewGetHandler(s)
	setHandler := handlers.NewSetHandler(s)
	pingHandler := handlers.NewPingHandler(s)
	echoHandler := handlers.NewEchoHandler(s)
	keysHandler := handlers.NewKeysHandler(s)
	infoHandler := handlers.NewInfoHandler(s)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection")
			continue
		}

		go func() {
			buff := make([]byte, 1024*8)
			for {
				conn.Read(buff)
				// query without value type mark
				err, cmd := command.ParseCommand(string(buff))
				if err != nil {
					log.Fatal("redis command error:", err)
				}

				switch (*cmd).Name() {
				case "CONFIG":
					configHandler.Handler(&conn, *cmd)
				case "GET":
					getHandler.Handler(&conn, *cmd)
				case "SET":
					setHandler.Handler(&conn, *cmd)
				case "ECHO":
					echoHandler.Handler(&conn, *cmd)
				case "PING":
					pingHandler.Handler(&conn, *cmd)
				case "KEYS":
					keysHandler.Handler(&conn, *cmd)
				case "INFO":
					infoHandler.Handler(&conn, *cmd)
				}
			}
		}()
	}
}
