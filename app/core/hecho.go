package core

func handleEcho(_ RedisInstance, conn Conn, args []string) {
	RespondString(conn, ToRedisSimpleString(args[1]))
}
