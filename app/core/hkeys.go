package core

func handleKeys(h RedisInstance, conn Conn, args []string) {
	keys := h.GetStore().GetKeys(args[1])
	RespondString(conn, ToRedisStrings(keys))
}
