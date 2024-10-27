package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/repr"
	"log"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/commands"
	"github.com/codecrafters-io/redis-starter-go/app/redis"
)

func NewEchoHandler(store *redis.Store) *EchoHandler {
	return &EchoHandler{
		store: *store,
	}
}

type EchoHandler struct {
	store redis.Store
}

func (h *EchoHandler) Handler(conn *net.Conn, c commands.Command[string]) {
	if c == nil || !c.Validate() {
		log.Fatal()
	}

	(*conn).Write(repr.ToRegularString(c.Args()[1]))
}
