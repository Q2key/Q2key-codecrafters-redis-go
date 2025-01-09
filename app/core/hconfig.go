package core

func handleConfig(h RedisInstance, conn Conn, args []string) {
	action, key := args[1], args[2]

	if action == "GET" && key == "dir" {
		resp := []string{key, h.GetConfig().Dir}
		RespondString(conn, ToRedisStrings(resp))
		return
	}

	if action == "GET" && key == "dbfilename" {
		resp := []string{key, h.GetConfig().DbFileName}
		RespondString(conn, ToRedisStrings(resp))
		return
	}

	RespondString(conn, ToRedisNullBulkString())
}
