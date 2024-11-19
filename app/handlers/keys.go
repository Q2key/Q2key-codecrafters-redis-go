package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/adapters"
	"github.com/codecrafters-io/redis-starter-go/app/contracts"
	"log"
	"net"
)

func NewKeysHandler(instance contracts.Instance) *KeysHandler {
	return &KeysHandler{
		instance: instance,
	}
}

type KeysHandler struct {
	instance contracts.Instance
}

func (h *KeysHandler) Handle(conn *net.Conn, c contracts.Command[string]) {
	if c == nil || !c.Validate() {
		log.Fatal()
	}

	args := c.Args()

	t := args[1]
	keys := h.instance.GetKeys(t)

	(*conn).Write([]byte(adapters.FromStringsArray(keys)))
}
