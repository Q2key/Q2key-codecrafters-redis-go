package contracts

import "net"

type Instance interface {
	Get(string) Value
	GetConfig() Config
	GetKeys(string) []string
	GetReplicaId() string
	GetMasterConn() *net.Conn
	GetChan() chan bool
	Set(key string, value string)
	GetStore() *map[string]Value
	SetExpiredAt(string, uint64)
	SetExpiredIn(string, uint64)
	RegisterReplicaConn(conn net.Conn)
	RegisterMasterConn(conn net.Conn)
	Propagate([]byte)
	Replicate([]byte)
}
