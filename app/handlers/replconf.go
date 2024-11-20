package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/contracts"
	"github.com/codecrafters-io/redis-starter-go/app/core"
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

func (h *ReplConfHandler) Handle(conn net.Conn, c contracts.Command) {
	if c == nil || !c.Validate() {
		log.Fatal()
	}

	conn.Write([]byte(core.FromStringToRedisCommonString("OK")))
}
