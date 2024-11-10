package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/command"
	"github.com/codecrafters-io/redis-starter-go/app/core"
	"github.com/codecrafters-io/redis-starter-go/app/repr"
	"log"
	"net"
)

func NewEchoHandler(instance *core.Instance) *EchoHandler {
	return &EchoHandler{
		instance: instance,
	}
}

type EchoHandler struct {
	instance *core.Instance
}

func (h *EchoHandler) Handler(conn *net.Conn, c command.Command[string]) {
	if c == nil || !c.Validate() {
		log.Fatal()
	}

	(*conn).Write([]byte(repr.FromString(c.Args()[1])))
}
