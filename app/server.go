package main

import (
	"bufio"
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
		buff := make([]byte, 1024)

		go func() {
			if ins.GetConfig().IsMaster() == false && (ins) != nil && (*ins.GetMasterConn()) != nil {
				for {
					res := make([]byte, 512)
					rdr := bufio.NewReader(*ins.GetMasterConn())
					n, _ := rdr.Read(res)

					str := string(res[:n])
					repc := ""

					for i, ch := range str {
						if string(ch) == "*" {
							repc = str[i:]
							break
						}
					}

					_, arr := core.FromRedisStringToStringArray(repc)
					for i, v := range arr {
						if v == "SET" && i+2 <= len(arr) {
							ins.Set(arr[i+1], arr[i+2])
						}
					}
				}
			}
		}()

		go func() {
			for {
				n, err := conn.Read(buff)

				if err == io.EOF {
					fmt.Println("Connection closed by client.")
					break
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

				if cmd.IsWrite() && ins.GetConfig().IsMaster() {
					repData := core.FromStringArrayToRedisStringArray(cmd.Args())
					ins.Replicate([]byte(repData))
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
