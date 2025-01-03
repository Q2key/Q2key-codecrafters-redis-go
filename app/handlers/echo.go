package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/contracts"
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

func (h *EchoHandler) Handle(conn contracts.Connection, args []string, _ *[]byte) {
	conn.Conn().Write([]byte(core.FromStringToRedisCommonString(args[1])))
}
