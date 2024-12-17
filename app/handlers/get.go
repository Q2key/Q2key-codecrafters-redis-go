package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/core"
	"github.com/codecrafters-io/redis-starter-go/app/core/rconn"
	"github.com/codecrafters-io/redis-starter-go/app/core/repr"
)

func NewGetHandler(instance core.Redis) *GetHandler {
	return &GetHandler{
		instance: instance,
	}
}

type GetHandler struct {
	instance core.Redis
}

func (h *GetHandler) Handle(conn rconn.RConn, args []string) {
	key := args[1]
	val, _ := h.instance.Store.Get(key)
	if val.IsExpired() {
		conn.Conn.Write([]byte(repr.ToRedisErrorString()))
	} else {
		conn.Conn.Write([]byte(repr.FromStringToRedisCommonString(val.GetValue())))
	}
}
