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

func (h *TypeHandler) Handle(conn rconn.RConn, args []string, _ *[]byte) {
	if len(args) < 1 {
		return
	}

	key := args[1]

	val, ok := h.instance.Store.Get(key)
	if !ok {
		conn.Conn.Write([]byte(repr.FromStringToRedisCommonString("none")))
	} else {
		conn.Conn.Write([]byte(repr.FromStringToRedisCommonString(val.ValueType)))
	}
}
