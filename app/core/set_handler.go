package core

import (
	"strconv"
)

func NewSetHandler(instance Redis) *SetHandler {
	return &SetHandler{
		instance: instance,
	}
}

type SetHandler struct {
	instance Redis
}

func (h *SetHandler) Handle(conn Conn, args []string, raw *[]byte) {
	key, val := args[1], args[2]

	vtype := GetValueTypes(string(*raw))

	_, ok := vtype[key]

	if !ok {
		return
	}

	h.instance.Store.Set(key, val)

	if len(args) >= 4 {
		exp, _ := strconv.Atoi(args[4])
		h.instance.Store.SetExpiredIn(key, uint64(exp))
	}

	conn.Conn().Write([]byte(FromStringToRedisCommonString("OK")))
}
