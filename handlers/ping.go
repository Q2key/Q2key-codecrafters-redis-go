package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/repr"
	"log"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/commands"
	"github.com/codecrafters-io/redis-starter-go/app/redis"
)

func NewPingHandler(store *redis.Store) *PingHandler {
	return &PingHandler{
		store: *store,
	}
}

type PingHandler struct {
	store redis.Store
}

func (h *PingHandler) Handler(conn *net.Conn, c commands.Command[string]) {
	if c == nil || !c.Validate() {
		log.Fatal()
	}

	(*conn).Write(repr.ToRegularString("PONG"))
}
