package core

func handleInfo(h RedisInstance, conn Conn, args []string) {
	r := h.GetConfig().GetReplica()
	res := "role:master"
	if r != nil {
		res = "role:slave"
	}

	for _, v := range args {
		if v == "replication" {
			res += ":master_replid:8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb"
			res += ":master_repl_offset:0"
		}
	}

	RespondString(conn, ToRedisBulkString(res))
}
