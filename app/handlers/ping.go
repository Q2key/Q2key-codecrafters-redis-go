package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/contracts"
	"github.com/codecrafters-io/redis-starter-go/app/core"
)

func NewPingHandler(instance core.Redis) *PingHandler {
	return &PingHandler{
		instance: instance,
	}
}

type PingHandler struct {
	instance core.Redis
}

func (h *PingHandler) Handle(conn contracts.Connector, _ []string, _ *[]byte) {
	conn.Conn().Write([]byte(core.FromStringToRedisCommonString("PONG")))
}
