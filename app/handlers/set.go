package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/contracts"
	"github.com/codecrafters-io/redis-starter-go/app/core"
	"log"
	"strconv"
)

func NewSetHandler(instance contracts.Instance) *SetHandler {
	return &SetHandler{
		instance: instance,
	}
}

type SetHandler struct {
	instance contracts.Instance
}

func (h *SetHandler) Handle(conn contracts.RedisConn, c contracts.Command) {
	if c == nil || !c.Validate() {
		log.Fatal()
	}

	if !c.Validate() {
		log.Fatal()
	}

	args := c.Args()
	key, val := args[1], args[2]

	h.instance.GetStore().Set(key, val)

	if len(args) >= 4 {
		exp, _ := strconv.Atoi(args[4])
		h.instance.GetStore().SetExpiredIn(key, uint64(exp))
	}

	conn.Conn().Write([]byte(core.FromStringToRedisCommonString("OK")))
}
