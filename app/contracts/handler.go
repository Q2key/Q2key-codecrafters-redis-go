package contracts

type Handler interface {
	Handle(RedisConn, Command)
}
