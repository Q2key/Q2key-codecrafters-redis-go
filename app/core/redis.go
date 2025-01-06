package core

import (
	"net"
)

type RedisInstance interface {
	HandleTCP(conn net.Conn)
	Init()
	GetConfig() *Config
	GetStore() *Store
}
