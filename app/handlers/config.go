package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/contracts"
	"github.com/codecrafters-io/redis-starter-go/app/core"
	"log"
)

func NewConfigHandler(instance contracts.Instance) *ConfigHandler {
	return &ConfigHandler{
		instance: &instance,
	}
}

type ConfigHandler struct {
	instance *contracts.Instance
}

func (h *ConfigHandler) Handle(conn contracts.RedisConn, c contracts.Command) {
	if c == nil || !c.Validate() {
		log.Fatal()
	}

	args := c.Args()
	action, key := args[1], args[2]

	if action == "GET" && key == "dir" {
		resp := []string{key, (*h.instance).GetConfig().GetDir()}
		conn.GetConn().Write([]byte(core.StringsToRedisStrings(resp)))
		return
	}

	if action == "GET" && key == "dbfilename" {
		resp := []string{key, (*h.instance).GetConfig().GetDbFileName()}
		conn.GetConn().Write([]byte(core.StringsToRedisStrings(resp)))
		return
	}

	conn.GetConn().Write([]byte(core.ToRedisErrorString()))
}
