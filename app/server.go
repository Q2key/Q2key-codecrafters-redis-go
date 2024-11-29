package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

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
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			buff := make([]byte, 1024*8)
			for {
				n, err := conn.Read(buff)
				if err != nil {
					return
				}

				sli := buff[:n]
				err, cmd := commands.ParseCommand(string(sli))

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

				if ins.GetMasterConn() != nil && !ins.GetConfig().IsMaster() {
					b := bufio.NewReader(*ins.GetMasterConn())
					br := make([]byte, 1024)
					n, _ := b.Read(br)

					err, cmd := commands.ParseCommand(string(br[:n]))
					if err != nil {
						continue
					}

					args := cmd.Args()

					for i := 0; i < len(args); i++ {
						if args[i] == "SET" {
							wg.Add(1)
						}
					}

					for i, c := range args {
						if c == "SET" {
							ins.Set(args[i+1], args[i+2])
							wg.Done()
						}
					}

					fmt.Println(cmd.Args())

				}

				if cmd.IsWrite() && ins.GetConfig().IsMaster() {
					repData := core.FromStringArrayToRedisStringArray(cmd.Args())
					ins.Replicate([]byte(repData))
				}
			}
		}()
		wg.Wait()
	}
}

func main() {
	fmt.Println("Starting server...")
	config := core.NewConfig().FromArguments(os.Args)
	redis := core.NewRedisInstance(*config)
	RunInstance(redis)
}
