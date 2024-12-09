package contracts

import (
	"net"
	"time"
)

type Scheduler struct {
	WaitTill             time.Time `json:"_"`
	WaitTimeoutMS        int       `json:"activeTimeoutMS"`
	ActiveRepicasCount   int       `json:"activeRepicasCount"`
	PendingReplicasCount int       `json:"pendingReplicasCount"`
	TotalReplicasCount   int       `json:"totalReplicasCount"`
}

type Instance interface {
	Get(string) Value
	GetConfig() Config
	GetKeys(string) []string
	GetReplicaId() string
	GetMasterConn() *net.Conn
	GetReplicas() map[string]*net.Conn
	Set(key string, value string)
	GetStore() *map[string]Value
	SetExpiredAt(string, uint64)
	SetExpiredIn(string, uint64)
	RegisterReplicaConn(conn net.Conn)
	RegisterMasterConn(conn net.Conn)
	Replicate([]byte)
	HandShakeMaster(ch chan []byte)
	GetScheduler() *Scheduler
	ScheduleReplicas(int, int)
}
