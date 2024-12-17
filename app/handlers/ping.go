package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/core"
	"github.com/codecrafters-io/redis-starter-go/app/core/rconn"
	"github.com/codecrafters-io/redis-starter-go/app/core/repr"
)

func NewPingHandler(instance core.Redis) *PingHandler {
	return &PingHandler{
		instance: instance,
	}
}

type PingHandler struct {
	instance core.Redis
}

func (h *PingHandler) Handle(conn rconn.RConn, _ []string) {
	conn.Conn.Write([]byte(repr.FromStringToRedisCommonString("PONG")))
}
