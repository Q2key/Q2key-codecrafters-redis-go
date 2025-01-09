package core

func handlePing(_ RedisInstance, conn Conn, _ []string) {
	RespondString(conn, ToRedisSimpleString("PONG"))
}
