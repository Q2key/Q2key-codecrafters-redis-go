package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/command"
	"github.com/codecrafters-io/redis-starter-go/app/core"
	"github.com/codecrafters-io/redis-starter-go/app/repr"
	"log"
	"net"
)

func NewPingHandler(store *core.Instance) *PingHandler {
	return &PingHandler{
		store: *store,
	}
}

type PingHandler struct {
	store core.Instance
}

func (h *PingHandler) Handler(conn *net.Conn, c command.Command[string]) {
	if c == nil || !c.Validate() {
		log.Fatal()
	}

	(*conn).Write([]byte(repr.FromString("PONG")))
}
