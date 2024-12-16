package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/contracts"
	"github.com/codecrafters-io/redis-starter-go/app/core"
	"log"
)

func NewPingHandler(instance contracts.Instance) *PingHandler {
	return &PingHandler{
		instance: instance,
	}
}

type PingHandler struct {
	instance contracts.Instance
}

func (h *PingHandler) Handle(conn contracts.RedisConn, c contracts.Command) {
	if c == nil || !c.Validate() {
		log.Fatal()
	}

	conn.Conn().Write([]byte(core.FromStringToRedisCommonString("PONG")))
}
