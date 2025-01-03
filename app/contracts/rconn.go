package contracts

import "net"

type Connection interface {
	Conn() net.Conn
	Id() string
	Offset() int
	SetOffset(offset int)
}
