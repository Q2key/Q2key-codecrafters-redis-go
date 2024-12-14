package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/contracts"
	"github.com/codecrafters-io/redis-starter-go/app/core"
	"log"
)

func NewGetHandler(instance contracts.Instance) *GetHandler {
	return &GetHandler{
		instance: instance,
	}
}

type GetHandler struct {
	instance contracts.Instance
}

func (h *GetHandler) Handle(conn contracts.RedisConn, c contracts.Command) {
	if c == nil || !c.Validate() {
		log.Fatal()
	}

	key := c.Args()[1]
	val, _ := h.instance.GetStore().Get(key)
	if val == nil || val.IsExpired() {
		conn.GetConn().Write([]byte(core.ToRedisErrorString()))
	} else {
		conn.GetConn().Write([]byte(core.FromStringToRedisCommonString(val.GetValue())))
	}
}
