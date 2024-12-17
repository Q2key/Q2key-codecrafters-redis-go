package contracts

import "net"

type Connector interface {
	Conn() net.Conn
	Id() string
	Offset() int
	SetOffset(offset int)
}
