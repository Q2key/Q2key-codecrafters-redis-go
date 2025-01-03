package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/core"
)

func NewGetHandler(instance core.Redis) *GetHandler {
	return &GetHandler{
		instance: instance,
	}
}

type GetHandler struct {
	instance core.Redis
}

func (h *GetHandler) Handle(conn core.Conn, args []string, _ *[]byte) {
	key := args[1]
	val, _ := h.instance.Store.Get(key)
	if val == nil || val.IsExpired() {
		conn.Conn().Write([]byte(core.ToRedisErrorString()))
	} else {
		conn.Conn().Write([]byte(core.FromStringToRedisCommonString(val.GetValue())))
	}
}
