package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/command"
	"github.com/codecrafters-io/redis-starter-go/app/core"
	"log"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/repr"
)

func NewKeysHandler(instance *core.Instance) *KeysHandler {
	return &KeysHandler{
		instance: instance,
	}
}

type KeysHandler struct {
	instance *core.Instance
}

func (h *KeysHandler) Handler(conn *net.Conn, c command.Command[string]) {
	if c == nil || !c.Validate() {
		log.Fatal()
	}

	args := c.Args()

	t := args[1]
	keys := h.instance.Keys(t)

	(*conn).Write([]byte(repr.FromStringsArray(keys)))
}
