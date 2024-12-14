package main

import (
	"context"
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

// todo rename to RunMaster
// todo create run Slave or so...
func RunInstance(ctx context.Context, ins contracts.Instance) {
	port := ins.GetConfig().GetPort()

	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		log.Fatalf("\r\nFailed to bind to port %s", port)
	}

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
		"WAIT":     handlers.NewWaitHandler(ins),
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

func handleRedisConnection(conn net.Conn, ins contracts.Instance, handlers map[string]contracts.Handler) {
	defer conn.Close()
	buff := make([]byte, 215)

	isMaster := ins.GetConfig().IsMaster()

	for {
		n, err := conn.Read(buff)
		if err == io.EOF {
			break
		}

		_, cmd := commands.ParseCommand(string(buff[:n]))
		h := handlers[cmd.Name()]
		h.Handle(core.NewReplicMasterConn(&conn), cmd)

		if isMaster && cmd.IsWrite() {
			ins.SendToReplicas(buff[:n])
		}
	}
}

func main() {
	fmt.Println("Starting server...")
	ctx := context.Background()
	config := core.NewConfig().FromArguments(os.Args)
	redis := core.NewInstance(ctx, *config)
	RunInstance(ctx, redis)
}
