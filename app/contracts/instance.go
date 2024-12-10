package contracts

import (
	"fmt"
	"net"
	"time"
)

type Scheduler struct {
	WaitTill             time.Time `json:"WaitTill"`
	WaitTimeoutMS        int       `json:"activeTimeoutMS"`
	ReplicasTowait       int       `json:"replicasToWait"`
	ActiveRepicasCount   int       `json:"activeRepicasCount"`
	PendingReplicasCount int       `json:"pendingReplicasCount"`
	TotalReplicasCount   int       `json:"totalReplicasCount"`
}

func (r *Scheduler) IsPending() bool {
	return r.WaitTill.UnixNano() >= time.Now().UnixNano()
}

func (r *Scheduler) Suspend(amount, msec int) {
	fmt.Println("Suspend")
	r.WaitTimeoutMS = msec
	r.ReplicasTowait = amount
	r.PendingReplicasCount = amount
	r.ActiveRepicasCount = r.TotalReplicasCount - r.PendingReplicasCount
	r.WaitTill = time.Now().Add(time.Duration(msec) * time.Millisecond)
	fmt.Printf("\r\n%v", r)
}

func (r *Scheduler) Release() {
	fmt.Println("Release")
	r.PendingReplicasCount -= 1
	r.ActiveRepicasCount = r.TotalReplicasCount - r.PendingReplicasCount
}

func (r *Scheduler) Update() {
	fmt.Println("Update")
	r.ActiveRepicasCount = r.TotalReplicasCount - r.PendingReplicasCount
}

func (r *Scheduler) AddActiveReplica() {
	fmt.Println("AddActiveReplica")
	r.TotalReplicasCount += 1
	r.Update()
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
