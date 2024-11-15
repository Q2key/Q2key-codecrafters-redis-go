package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/contracts"
	"github.com/codecrafters-io/redis-starter-go/app/core"
	"github.com/codecrafters-io/redis-starter-go/app/repr"
	"log"
	"net"
)

func NewInfoHandler(store *core.Instance) *InfoHandler {
	return &InfoHandler{
		instance: store,
	}
}

type InfoHandler struct {
	instance *core.Instance
}

func (h *InfoHandler) Handler(conn *net.Conn, c contracts.Command[string]) {
	if c == nil || !c.Validate() {
		log.Fatal()
	}

	r := h.instance.Config.GetReplica()
	role := "role:master"
	if r != nil {
		role = "role:slave"
	}

	(*conn).Write([]byte(repr.BulkString(role)))
}
