package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/contracts"
	"github.com/codecrafters-io/redis-starter-go/app/repr"
	"log"
	"net"
)

func NewConfigHandler(instance contracts.Instance) *ConfigHandler {
	return &ConfigHandler{
		instance: &instance,
	}
}

type ConfigHandler struct {
	instance *contracts.Instance
}

func (h *ConfigHandler) Handle(conn *net.Conn, c contracts.Command[string]) {
	if c == nil || !c.Validate() {
		log.Fatal()
	}

	args := c.Args()
	action, key := args[1], args[2]

	if action == "GET" && key == "dir" {
		resp := []string{key, (*h.instance).GetConfig().GetDir()}
		(*conn).Write([]byte(repr.FromStringsArray(resp)))
		return
	}

	if action == "GET" && key == "dbfilename" {
		resp := []string{key, (*h.instance).GetConfig().GetDbFileName()}
		(*conn).Write([]byte(repr.FromStringsArray(resp)))
		return
	}

	(*conn).Write([]byte(repr.ErrorString()))
}
