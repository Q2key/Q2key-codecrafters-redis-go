package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/core"
)

func NewKeysHandler(instance core.Redis) *KeysHandler {
	return &KeysHandler{
		instance: instance,
	}
}

type KeysHandler struct {
	instance core.Redis
}

func (h *KeysHandler) Handle(conn core.Conn, args []string, _ *[]byte) {
	t := args[1]
	keys := h.instance.Store.GetKeys(t)

	conn.Conn().Write([]byte(core.StringsToRedisStrings(keys)))
}
