package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/contracts"
	"github.com/codecrafters-io/redis-starter-go/app/core"
)

func NewConfigHandler(instance core.Redis) *ConfigHandler {
	return &ConfigHandler{
		instance: &instance,
	}
}

type ConfigHandler struct {
	instance *core.Redis
}

func (h *ConfigHandler) Handle(conn contracts.Connector, args []string, _ *[]byte) {
	action, key := args[1], args[2]

	if action == "GET" && key == "dir" {
		resp := []string{key, (*h.instance).Config.Dir}
		conn.Conn().Write([]byte(core.StringsToRedisStrings(resp)))
		return
	}

	if action == "GET" && key == "dbfilename" {
		resp := []string{key, (*h.instance).Config.DbFileName}
		conn.Conn().Write([]byte(core.StringsToRedisStrings(resp)))
		return
	}

	conn.Conn().Write([]byte(core.ToRedisErrorString()))
}
