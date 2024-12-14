package core

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/contracts"
)

type RedisConn struct {
	Conn      net.Conn
	id        string
	BytesSent int
}

func (r *RedisConn) GetId() string {
	return r.id
}

func (r *RedisConn) GetConn() net.Conn {
	return r.Conn
}

func (r *RedisConn) GetOffset() int {
	return r.BytesSent
}

func (r *RedisConn) SetOffset(offset int) {
	r.BytesSent = offset
}

func NewReplicaConn(conn *net.Conn, replicaId string) contracts.RedisConn {
	return &RedisConn{
		Conn:      *conn,
		id:        replicaId,
		BytesSent: 0,
	}
}

func NewReplicMasterConn(conn *net.Conn) contracts.RedisConn {
	return &RedisConn{
		Conn:      *conn,
		id:        RandStringBytes(10),
		BytesSent: 0,
	}
}
