package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/adapters"
	"github.com/codecrafters-io/redis-starter-go/app/contracts"
	"log"
	"net"
)

func NewReplConfHandler(instance contracts.Instance) *ReplConfHandler {
	return &ReplConfHandler{
		instance: instance,
	}
}

type ReplConfHandler struct {
	instance contracts.Instance
}

func (h *ReplConfHandler) Handle(conn *net.Conn, c contracts.Command[string]) {
	if c == nil || !c.Validate() {
		log.Fatal()
	}

	(*conn).Write([]byte(adapters.FromString("OK")))
}
