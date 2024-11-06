package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/command"
	"github.com/codecrafters-io/redis-starter-go/app/core"
	"github.com/codecrafters-io/redis-starter-go/app/repr"
	"log"
	"net"
)

func NewGetHandler(store *core.Instance) *GetHandler {
	return &GetHandler{
		store: store,
	}
}

type GetHandler struct {
	store *core.Instance
}

func (h *GetHandler) Handler(conn *net.Conn, c command.Command[string]) {
	if c == nil || !c.Validate() {
		log.Fatal()
	}

	key := c.Args()[1]
	val := h.store.Get(key)

	if val.IsExpired() {
		(*conn).Write([]byte(repr.ErrorString()))
	} else {
		(*conn).Write([]byte(repr.FromString(val.Val)))
	}
}
