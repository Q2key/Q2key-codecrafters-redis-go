package handlers

import (
	"log"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/commands"
	"github.com/codecrafters-io/redis-starter-go/app/redis"
	"github.com/codecrafters-io/redis-starter-go/app/repr"
)

func NewKeysHandler(store *redis.Store) *KeysHandler {
	return &KeysHandler{
		store: *store,
	}
}

type KeysHandler struct {
	store redis.Store
}

func (h *KeysHandler) Handler(conn *net.Conn, c commands.Command[string]) {
	if c == nil || !c.Validate() {
		log.Fatal()
	}

	(*conn).Write(repr.ToErrorString(nil))
}
