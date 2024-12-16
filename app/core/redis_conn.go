package core

import (
	"github.com/codecrafters-io/redis-starter-go/app/contracts"
	"net"
)

type RedisConn struct {
	conn   net.Conn
	id     string
	offset int
}

func (r *RedisConn) Id() string {
	return r.id
}

func (r *RedisConn) Conn() net.Conn {
	return r.conn
}

func (r *RedisConn) Offset() int {
	return r.offset
}

func (r *RedisConn) SetOffset(offset int) {
	r.offset = offset
}

func NewRedisConn(conn *net.Conn) contracts.RedisConn {
	return &RedisConn{
		conn:   *conn,
		id:     RandStringBytes(10),
		offset: 0,
	}
}

func NewRedisConcreteConn(conn *net.Conn) *RedisConn {
	return &RedisConn{
		conn:   *conn,
		id:     RandStringBytes(10),
		offset: 0,
	}
}
