package handlers

import (
	"log"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/contracts"
	"github.com/codecrafters-io/redis-starter-go/app/core"
)

func NewGetHandler(instance contracts.Instance) *GetHandler {
	return &GetHandler{
		instance: instance,
	}
}

type GetHandler struct {
	instance contracts.Instance
}

func (h *GetHandler) Handle(conn net.Conn, c contracts.Command) {
	if c == nil || !c.Validate() {
		log.Fatal()
	}

	key := c.Args()[1]
	val := h.instance.Get(key)

	if val == nil || val.IsExpired() {
		conn.Write([]byte(core.ToRedisErrorString()))
	} else {
		conn.Write([]byte(core.FromStringToRedisCommonString(val.GetValue())))
	}
}
