package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/core"
)

func RunInstance(ctx context.Context, ins core.RedisInstance) {
	port := ins.GetConfig().Port

	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		log.Fatalf("\r\nFailed to bind to port %s", port)
	}

	ins.Init()
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
	var instance core.RedisInstance
	if config.IsMaster() {
		instance = core.NewMaster(ctx, *config)
	} else {
		instance = core.NewReplica(ctx, *config)
	}
	RunInstance(ctx, instance)
}
