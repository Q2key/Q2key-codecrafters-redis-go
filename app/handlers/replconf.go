package handlers

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/core/rconn"
	"github.com/codecrafters-io/redis-starter-go/app/core/repr"

	"github.com/codecrafters-io/redis-starter-go/app/core"
)

func NewReplConfHandler(instance core.Redis) *ReplConfHandler {
	return &ReplConfHandler{
		instance: instance,
	}
}

type ReplConfHandler struct {
	instance core.Redis
}

func (h *ReplConfHandler) Handle(conn rconn.RConn, args []string, _ *[]byte) {
	if len(args) > 2 && args[1] == "ACK" {
		cnt := args[2]
		num, _ := strconv.Atoi(cnt)
		id := conn.Id

		h.instance.UpdateReplica(id, num)
		*h.instance.AckChan <- rconn.Ack{ConnId: id, Offset: num}

		conn.Conn.Write([]byte(""))
	} else {
		conn.Conn.Write([]byte(repr.FromStringToRedisCommonString("OK")))
	}
}
