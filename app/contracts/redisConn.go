package contracts

import "net"

type RedisConn interface {
	GetId() string
	GetConn() net.Conn
	GetOffset() int
	SetOffset(offset int)
}
