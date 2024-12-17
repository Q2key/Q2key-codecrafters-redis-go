package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/core"
	"github.com/codecrafters-io/redis-starter-go/app/core/rconn"
	"github.com/codecrafters-io/redis-starter-go/app/core/repr"
)

func NewKeysHandler(instance core.Redis) *KeysHandler {
	return &KeysHandler{
		instance: instance,
	}
}

type KeysHandler struct {
	instance core.Redis
}

func (h *KeysHandler) Handle(conn rconn.RConn, args []string) {
	t := args[1]
	keys := h.instance.Store.GetKeys(t)

	conn.Conn.Write([]byte(repr.StringsToRedisStrings(keys)))
}
