package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/command"
	"github.com/codecrafters-io/redis-starter-go/app/core"
	"github.com/codecrafters-io/redis-starter-go/app/repr"
	"log"
	"net"
)

func NewEchoHandler(store *core.Instance) *EchoHandler {
	return &EchoHandler{
		store: store,
	}
}

type EchoHandler struct {
	store *core.Instance
}

func (h *EchoHandler) Handler(conn *net.Conn, c command.Command[string]) {
	if c == nil || !c.Validate() {
		log.Fatal()
	}

	(*conn).Write([]byte(repr.FromString(c.Args()[1])))
}
