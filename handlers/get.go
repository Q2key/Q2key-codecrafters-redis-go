package handlers

import (
	"log"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/commands"
	"github.com/codecrafters-io/redis-starter-go/app/redis"
	"github.com/codecrafters-io/redis-starter-go/app/repr"
)

func NewGetHandler(store *redis.Store) *GetHandler {
	return &GetHandler{
		store: *store,
	}
}

type GetHandler struct {
	store redis.Store
}

func (h *GetHandler) Handler(conn *net.Conn, c commands.Command[string]) {
	if c == nil || !c.Validate() {
		log.Fatal()
	}

	key := c.Args()[0]
	val := h.store.Get(key)
	if val.IsExpired() {
		(*conn).Write(repr.ToErrorString(nil))
	} else {
		(*conn).Write(repr.ToRegularString(val.Value))
	}
}
