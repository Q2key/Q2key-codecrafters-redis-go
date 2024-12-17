package handlers

import (
	"fmt"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/core"
	"github.com/codecrafters-io/redis-starter-go/app/core/rconn"
	"github.com/codecrafters-io/redis-starter-go/app/core/repr"
)

func NewSetHandler(instance core.Redis) *SetHandler {
	return &SetHandler{
		instance: instance,
	}
}

type SetHandler struct {
	instance core.Redis
}

func (h *SetHandler) Handle(conn rconn.RConn, args []string, raw *[]byte) {
	key, val := args[1], args[2]

	vtype := repr.GetValueTypes(string(*raw))

	v, ok := vtype[key]

	if !ok {
		return
	}

	fmt.Println(v)

	h.instance.Store.Set(key, val)

	if len(args) >= 4 {
		exp, _ := strconv.Atoi(args[4])
		h.instance.Store.SetExpiredIn(key, uint64(exp))
	}

	conn.Conn.Write([]byte(repr.FromStringToRedisCommonString("OK")))
}
