package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/contracts"
	"github.com/codecrafters-io/redis-starter-go/app/core"
	"log"
)

func NewTypeHandler(instance contracts.Instance) *TypeHandler {
	return &TypeHandler{
		instance: instance,
	}
}

type TypeHandler struct {
	instance contracts.Instance
}

func (h *TypeHandler) Handle(conn contracts.RedisConn, c contracts.Command) {
	if c == nil || !c.Validate() {
		log.Fatal()
	}

	conn.Conn().Write([]byte(core.FromStringToRedisCommonString("PONG")))
}
