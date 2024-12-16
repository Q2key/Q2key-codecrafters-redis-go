package contracts

import "net"

type RedisConn interface {
	Id() string
	Conn() net.Conn
	Offset() int
	SetOffset(offset int)
}
