package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/contracts"
	"github.com/codecrafters-io/redis-starter-go/app/core"
	"log"
)

func NewEchoHandler(instance contracts.Instance) *EchoHandler {
	return &EchoHandler{
		instance: instance,
	}
}

type EchoHandler struct {
	instance contracts.Instance
}

func (h *EchoHandler) Handle(conn contracts.RedisConn, c contracts.Command) {
	if c == nil || !c.Validate() {
		log.Fatal()
	}

	conn.GetConn().Write([]byte(core.FromStringToRedisCommonString(c.Args()[1])))
}
