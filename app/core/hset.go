package core

import "strconv"

func handleSet(h RedisInstance, conn Conn, args []string) {
	key, val := args[1], args[2]

	// according to the doc set is always setting the string value type
	h.GetStore().Set(key, val, STRING)

	if len(args) >= 4 {
		exp, _ := strconv.Atoi(args[4])
		h.GetStore().SetExpiredIn(key, uint64(exp))
	}

	RespondString(conn, ToRedisSimpleString("OK"))
}
