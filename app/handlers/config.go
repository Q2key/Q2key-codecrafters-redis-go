package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/command"
	"github.com/codecrafters-io/redis-starter-go/app/core"
	"github.com/codecrafters-io/redis-starter-go/app/repr"
	"log"
	"net"
)

func NewConfigHandler(store *core.Instance) *ConfigHandler {
	return &ConfigHandler{
		store: store,
	}
}

type ConfigHandler struct {
	store *core.Instance
}

func (h *ConfigHandler) Handler(conn *net.Conn, c command.Command[string]) {
	if c == nil || !c.Validate() {
		log.Fatal()
	}

	args := c.Args()
	action, key := args[1], args[2]

	if action == "GET" && key == "dir" {
		resp := []string{key, h.store.Config.GetDir()}
		(*conn).Write([]byte(repr.FromStringsArray(resp)))
		return
	}

	if action == "GET" && key == "dbfilename" {
		resp := []string{key, h.store.Config.GetDbFileName()}
		(*conn).Write([]byte(repr.FromStringsArray(resp)))
		return
	}

	(*conn).Write([]byte(repr.ErrorString()))
}
