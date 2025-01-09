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

type Redis struct {
	Store    Store
	Config   Config
	Commands map[string]CommandHandler
}
