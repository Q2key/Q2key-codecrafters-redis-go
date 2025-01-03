package core

import (
	"net"
)

type Ack struct {
	ConnId string
	Offset int
}

type Conn struct {
	conn   net.Conn
	id     string
	offset int
}

func NewRConn(conn *net.Conn) *Conn {
	return &Conn{
		conn:   *conn,
		id:     randStringBytes(10),
		offset: 0,
	}
}

func (r *Conn) Conn() net.Conn {
	return r.conn
}

func (r *Conn) Id() string {
	return r.id
}

func (r *Conn) Offset() int {
	return r.offset
}

func (r *Conn) SetOffset(offset int) {
	r.offset = offset
}
