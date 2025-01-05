package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/core"
)

func RunInstance(ctx context.Context, ins core.Redis) {
	port := ins.Config.Port

	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		log.Fatalf("\r\nFailed to bind to port %s", port)
	}

	go ins.Handshake()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("Something went wrong with tcp connection...")
		}

		go ins.HandleTCP(conn)
	}
}

func main() {
	ctx := context.Background()
	config := core.NewConfig().FromArguments(os.Args)
	redis := core.NewRedis(ctx, *config)
	RunInstance(ctx, *redis)
}
