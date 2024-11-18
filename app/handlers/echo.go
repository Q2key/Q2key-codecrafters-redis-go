package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/contracts"
	"github.com/codecrafters-io/redis-starter-go/app/mappers"
	"log"
	"net"
)

func NewEchoHandler(instance contracts.Instance) *EchoHandler {
	return &EchoHandler{
		instance: instance,
	}
}

type EchoHandler struct {
	instance contracts.Instance
}

func (h *EchoHandler) Handle(conn *net.Conn, c contracts.Command[string]) {
	if c == nil || !c.Validate() {
		log.Fatal()
	}

	(*conn).Write([]byte(mappers.FromString(c.Args()[1])))
}
