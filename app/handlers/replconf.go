package handlers

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/core"
)

func NewReplConfHandler(redis core.Redis) *ReplConfHandler {
	return &ReplConfHandler{
		redis: redis,
	}
}

type ReplConfHandler struct {
	redis core.Redis
}

func (h *ReplConfHandler) Handle(conn core.Conn, args []string, _ *[]byte) {
	if len(args) > 2 && args[1] == "ACK" {
		cnt := args[2]
		num, _ := strconv.Atoi(cnt)
		id := conn.Id()

		h.redis.UpdateReplica(id, num)
		*h.redis.AckChan <- core.Ack{ConnId: id, Offset: num}

		conn.Conn().Write([]byte(""))
	} else {
		conn.Conn().Write([]byte(core.FromStringToRedisCommonString("OK")))
	}
}
