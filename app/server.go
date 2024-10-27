package main

import (
	"fmt"
	"github.com/codecrafters-io/redis-starter-go/app/repr"
	"github.com/codecrafters-io/redis-starter-go/handlers"
	"log"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/redis"
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

	defer ln.Close()

	args := os.Args
	cfg := &redis.Config{}
	if len(args) > 1 {
		for i := 1; i < len(args); i++ {
			a := args[i]
			if i+1 == len(args) {
				break
			}

			v := args[i+1]
			if a == "--dir" {
				cfg.SetDir(v)
			}

			if a == "--dbfilename" {
				cfg.SetDbFilenabe(v)
			}
		}
	}

	store := redis.NewRedisStore()

	configHandler := handlers.NewConfigHandler(&store)
	getHandler := handlers.NewGetHandler(&store)
	setHandler := handlers.NewSetHandler(&store)
	pingHandler := handlers.NewPingHandler(&store)
	echoHandler := handlers.NewEchoHandler(&store)
	keysHandler := handlers.NewKeysHandler(&store)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection")
			continue
		}

		go func() {
			buff := make([]byte, 2024)
			for {
				conn.Read(buff)
				// query without value type mark
				err, cmd := repr.ParseCommand(string(buff))
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
				}
			}
		}()
	}
}
