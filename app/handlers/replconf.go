package handlers

import (
	"log"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/contracts"
	"github.com/codecrafters-io/redis-starter-go/app/core"
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

	if len(c.Args()) > 2 && c.Args()[1] == "ACK" {
		s := h.instance.GetScheduler()
		s.Release()
	}

	conn.Write([]byte(core.FromStringToRedisCommonString("OK")))
}
