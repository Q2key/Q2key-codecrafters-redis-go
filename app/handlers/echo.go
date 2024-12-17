package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/core"
	"github.com/codecrafters-io/redis-starter-go/app/core/rconn"
	"github.com/codecrafters-io/redis-starter-go/app/core/repr"
)

func NewEchoHandler(instance core.Redis) *EchoHandler {
	return &EchoHandler{
		instance: instance,
	}
}

type EchoHandler struct {
	instance core.Redis
}

func (h *EchoHandler) Handle(conn rconn.RConn, args []string) {
	conn.Conn.Write([]byte(repr.FromStringToRedisCommonString(args[1])))
}
