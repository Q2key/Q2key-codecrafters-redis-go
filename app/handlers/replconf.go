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

	if len(c.Args()) > 2 && c.Args()[1] == "ACK" {
		cnt := c.Args()[2]
		num, _ := strconv.Atoi(cnt)
		id := conn.Id()

		h.instance.UpdateReplica(id, num)
		*h.instance.GetAckChan() <- contracts.Ack{ConnId: id, Offset: num}

		conn.Conn().Write([]byte(""))
	} else {
		conn.Conn().Write([]byte(core.FromStringToRedisCommonString("OK")))
	}
}
