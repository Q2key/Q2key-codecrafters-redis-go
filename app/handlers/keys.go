package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/contracts"
	"github.com/codecrafters-io/redis-starter-go/app/core"
	"log"
)

func NewKeysHandler(instance contracts.Instance) *KeysHandler {
	return &KeysHandler{
		instance: instance,
	}
}

type KeysHandler struct {
	instance contracts.Instance
}

func (h *KeysHandler) Handle(conn contracts.RedisConn, c contracts.Command) {
	if c == nil || !c.Validate() {
		log.Fatal()
	}

	args := c.Args()

	t := args[1]
	keys := h.instance.GetStore().GetKeys(t)

	conn.Conn().Write([]byte(core.StringsToRedisStrings(keys)))
}
