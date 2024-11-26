package contracts

import "net"

type Instance interface {
	Get(string) Value
	GetConfig() Config
	GetKeys(string) []string
	GetReplicaId() string
	Set(key string, value string)
	GetStore() *map[string]Value
	SetExpiredAt(string, uint64)
	SetExpiredIn(string, uint64)
	SetRemoteAddr(string)
	GetRemoteAddr() string
	SetReplicaConn(conn net.Conn)
	GetReplicaConn() *net.Conn
	Replicate([]byte)
}
