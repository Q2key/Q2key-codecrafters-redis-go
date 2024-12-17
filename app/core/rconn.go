package core

import (
	"net"
)

type Ack struct {
	ConnId string
	Offset int
}

type RConn struct {
	Conn   net.Conn
	Id     string
	Offset int
}

func NewRConn(conn *net.Conn) *RConn {
	return &RConn{
		Conn:   *conn,
		Id:     randStringBytes(10),
		Offset: 0,
	}
}
