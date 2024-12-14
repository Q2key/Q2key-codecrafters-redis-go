package handlers

import (
	"log"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/contracts"
	"github.com/codecrafters-io/redis-starter-go/app/core"
)

func NewReplConfHandler(instance contracts.Instance) *ReplConfHandler {
	return &ReplConfHandler{
		instance: instance,
	}
}

type ReplConfHandler struct {
	instance contracts.Instance
}

func (h *ReplConfHandler) Handle(conn contracts.RedisConn, c contracts.Command) {
	if c == nil || !c.Validate() {
		log.Fatal()
	}

	if len(c.Args()) > 2 && c.Args()[1] == "listening-port" {
		h.instance.RegisterReplicaConn(&conn)
	}

	if len(c.Args()) > 2 && c.Args()[1] == "ACK" {
		cnt := c.Args()[2]
		num, _ := strconv.Atoi(cnt)
		*h.instance.GetAckChan() <- contracts.Ack{ConnId: conn.GetId(), Offset: num}
	}

	conn.GetConn().Write([]byte(core.FromStringToRedisCommonString("OK")))
}
