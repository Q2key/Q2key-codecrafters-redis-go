package core

import "strconv"

func handleReplconf(r RedisInstance, conn Conn, args []string) {
	master, ok := r.(*Master)
	if !ok {
		return
	}
	h := master
	if len(args) > 2 && args[1] == "ACK" {
		cnt := args[2]
		num, _ := strconv.Atoi(cnt)
		id := conn.Id()

		h.UpdateReplica(id, num)
		*h.AckChan <- Ack{ConnId: id, Offset: num}
	} else {
		RespondString(conn, ToRedisSimpleString("OK"))
	}
}
