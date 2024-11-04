package handlers

import (
	"fmt"
	"github.com/codecrafters-io/redis-starter-go/app/command"
	"github.com/codecrafters-io/redis-starter-go/app/core"
	"github.com/codecrafters-io/redis-starter-go/app/repr"
	"log"
	"net"
	"strconv"
)

func NewSetHandler(store *core.Instance) *SetHandler {
	return &SetHandler{
		store: store,
	}
}

type SetHandler struct {
	store *core.Instance
}

func (h *SetHandler) Handler(conn *net.Conn, c command.Command[string]) {
	fmt.Printf("\r\ninvoke set handler %v", c.Args())
	if c == nil || !c.Validate() {
		log.Fatal()
	}

	if !c.Validate() {
		log.Fatal()
	}

	args := c.Args()
	key, val := args[1], args[2]
	var exp int

	if len(args) >= 4 {
		exp, _ = strconv.Atoi(args[4])
	}

	h.store.Set(key, val, int64(exp))

	fmt.Printf("\r\nstore %v", h.store.Store())

	(*conn).Write([]byte(repr.FromString("OK")))
}
