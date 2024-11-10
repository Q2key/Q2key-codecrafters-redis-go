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
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	ln, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	//defer ln.Close()

	s := core.NewRedisInstanceWithArgs(os.Args)

	configHandler := handlers.NewConfigHandler(s)
	getHandler := handlers.NewGetHandler(s)
	setHandler := handlers.NewSetHandler(s)
	pingHandler := handlers.NewPingHandler(s)
	echoHandler := handlers.NewEchoHandler(s)
	keysHandler := handlers.NewKeysHandler(s)

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

				fmt.Printf("[%v]", (*cmd).Args())

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
				}
			}
		}()
	}
}
