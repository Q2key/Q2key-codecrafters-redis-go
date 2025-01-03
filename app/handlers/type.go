package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/contracts"
	"github.com/codecrafters-io/redis-starter-go/app/core"
)

func NewTypeHandler(instance core.Redis) *TypeHandler {
	return &TypeHandler{
		instance: instance,
	}
}

type TypeHandler struct {
	instance core.Redis
}

func (h *TypeHandler) Handle(conn contracts.Connector, args []string, _ *[]byte) {
	if len(args) < 1 {
		return
	}

	key := args[1]

	val, ok := h.instance.Store.Get(key)
	if !ok {
		conn.Conn().Write([]byte(core.FromStringToRedisCommonString("none")))
	} else {
		conn.Conn().Write([]byte(core.FromStringToRedisCommonString(val.ValueType)))
	}
}
