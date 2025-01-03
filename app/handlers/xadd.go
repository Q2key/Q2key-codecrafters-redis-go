package handlers

import "github.com/codecrafters-io/redis-starter-go/app/core"

func NewXaddHandler(instance core.Redis) *XaddHandler {
	return &XaddHandler{
		instance: instance,
	}
}

type XaddHandler struct {
	instance core.Redis
}

func (h *XaddHandler) Handle(conn core.Conn, args []string, _ *[]byte) {
	if len(args) < 1 {
		return
	}

	key := args[1]

	conn.Conn().Write([]byte(core.FromStringToRedisCommonString(key)))
}
