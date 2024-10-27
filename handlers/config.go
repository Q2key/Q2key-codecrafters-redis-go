package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/redis"
	"log"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/commands"
	"github.com/codecrafters-io/redis-starter-go/app/repr"
)

func NewConfigHandler(store *redis.Store) *ConfigHandler {
	return &ConfigHandler{
		store: *store,
	}
}

type ConfigHandler struct {
	store redis.Store
}

func (h *ConfigHandler) Handler(conn *net.Conn, c commands.Command[string]) {
	if c == nil || !c.Validate() {
		log.Fatal()
	}

	args := c.Args()
	action, key := args[1], args[2]

	if action == "GET" && key == "dir" {
		resp := []string{key, h.store.GetConfig().GetDir()}
		(*conn).Write(repr.ToStringArray(resp))
		return
	}

	if action == "GET" && key == "dbfilename" {
		resp := []string{key, h.store.GetConfig().GetDbFileName()}
		(*conn).Write(repr.ToStringArray(resp))
		return
	}

	(*conn).Write(repr.ToErrorString(nil))
}
