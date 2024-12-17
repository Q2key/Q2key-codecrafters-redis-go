package handlers

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/core"
)

func NewSetHandler(instance core.Redis) *SetHandler {
	return &SetHandler{
		instance: instance,
	}
}

type SetHandler struct {
	instance core.Redis
}

func (h *SetHandler) Handle(conn core.RConn, args []string, raw *[]byte) {
	key, val := args[1], args[2]

	vtype := core.GetValueTypes(string(*raw))

	_, ok := vtype[key]

	if !ok {
		return
	}

	h.instance.Store.Set(key, val)

	if len(args) >= 4 {
		exp, _ := strconv.Atoi(args[4])
		h.instance.Store.SetExpiredIn(key, uint64(exp))
	}

	conn.Conn.Write([]byte(core.FromStringToRedisCommonString("OK")))
}
