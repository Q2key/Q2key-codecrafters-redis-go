package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/contracts"
	"github.com/codecrafters-io/redis-starter-go/app/core"
	"github.com/codecrafters-io/redis-starter-go/app/repr"
	"log"
	"net"
	"strconv"
)

func NewSetHandler(store *core.Instance) *SetHandler {
	return &SetHandler{
		instance: store,
	}
}

type SetHandler struct {
	instance *core.Instance
}

func (h *SetHandler) Handler(conn *net.Conn, c contracts.Command[string]) {
	if c == nil || !c.Validate() {
		log.Fatal()
	}

	if !c.Validate() {
		log.Fatal()
	}

	args := c.Args()
	key, val := args[1], args[2]

	h.instance.Set(key, val)

	if len(args) >= 4 {
		exp, _ := strconv.Atoi(args[4])
		h.instance.SetExpiredIn(key, uint64(exp))
	}

	(*conn).Write([]byte(repr.FromString("OK")))
}
