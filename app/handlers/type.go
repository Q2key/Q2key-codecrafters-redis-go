package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/core"
	"github.com/codecrafters-io/redis-starter-go/app/core/rconn"
	"github.com/codecrafters-io/redis-starter-go/app/core/repr"
)

func NewTypeHandler(instance core.Redis) *TypeHandler {
	return &TypeHandler{
		instance: instance,
	}
}

type TypeHandler struct {
	instance core.Redis
}

func (h *TypeHandler) Handle(conn rconn.RConn, _ []string) {
	conn.Conn.Write([]byte(repr.FromStringToRedisCommonString("PONG")))
}
