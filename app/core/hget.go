package core

func handleGet(h RedisInstance, conn Conn, args []string) {
	key := args[1]
	val, _ := h.GetStore().Get(key)
	if val == nil || val.IsExpired() {
		RespondString(conn, ToRedisNullBulkString())
	} else {
		RespondString(conn, ToRedisSimpleString(val.ToString()))
	}
}
