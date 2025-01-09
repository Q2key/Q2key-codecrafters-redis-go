package core

func handleType(h RedisInstance, conn Conn, args []string) {
	if len(args) < 1 {
		return
	}

	key := args[1]

	val, ok := h.GetStore().Get(key)
	if !ok {
		RespondString(conn, ToRedisSimpleString("none"))
	} else {
		RespondString(conn, ToRedisSimpleString(string(val.GetType())))
	}
}
