package core

import (
	"strconv"
)

func NewReplConfHandler(redis Redis) *ReplConfHandler {
	return &ReplConfHandler{
		redis: redis,
	}
}

type ReplConfHandler struct {
	redis Redis
}

func (h *ReplConfHandler) Handle(conn Conn, args []string, _ *[]byte) {
	if len(args) > 2 && args[1] == "ACK" {
		cnt := args[2]
		num, _ := strconv.Atoi(cnt)
		id := conn.Id()

		h.redis.UpdateReplica(id, num)
		*h.redis.AckChan <- Ack{ConnId: id, Offset: num}

		conn.Conn().Write([]byte(""))
	} else {
		conn.Conn().Write([]byte(FromStringToRedisCommonString("OK")))
	}
}
