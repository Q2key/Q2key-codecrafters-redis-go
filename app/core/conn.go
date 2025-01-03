package core

import "net"

type Conn interface {
	Conn() net.Conn
	Id() string
	Offset() int
	SetOffset(offset int)
}
