package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/core"
)

func NewEchoHandler(instance core.Redis) *EchoHandler {
	return &EchoHandler{
		instance: instance,
	}
}

type EchoHandler struct {
	instance core.Redis
}

func (h *EchoHandler) Handle(conn core.RConn, args []string, _ *[]byte) {
	conn.Conn.Write([]byte(core.FromStringToRedisCommonString(args[1])))
}
