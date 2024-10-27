package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/repr"
	"log"
	"net"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/commands"
	"github.com/codecrafters-io/redis-starter-go/app/redis"
)

func NewSetHandler(store *redis.Store) *SetHandler {
	return &SetHandler{
		store: *store,
	}
}

type SetHandler struct {
	store redis.Store
}

func (h *SetHandler) Handler(conn *net.Conn, c commands.Command[string]) {
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
	(*conn).Write(repr.ToRegularString("OK"))
}
