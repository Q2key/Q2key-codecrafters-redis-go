package core

import (
	"github.com/codecrafters-io/redis-starter-go/app/contracts"
	"net"
)

type Ack struct {
	ConnId string
	Offset int
}

type RConn struct {
	conn   net.Conn
	id     string
	offset int
}

func NewRConn(conn *net.Conn) contracts.Connector {
	return &RConn{
		conn:   *conn,
		id:     randStringBytes(10),
		offset: 0,
	}
}

func (r *RConn) Conn() net.Conn {
	return r.conn
}

func (r *RConn) Id() string {
	return r.id
}

func (r *RConn) Offset() int {
	return r.offset
}

func (r *RConn) SetOffset(offset int) {
	r.offset = offset
}
