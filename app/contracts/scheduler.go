package contracts

import (
	"encoding/json"
	"fmt"
	"time"
)

type Scheduler struct {
	WaitTill                  *time.Time `json:"WaitTill"`
	WaitTimeoutMS             *int       `json:"activeTimeoutMS"`
	ReplicasToWait            int        `json:"replicasToWait"`
	TotalReplicasCount        int        `json:"totalReplicasCount"`
	AcknowledgedReplicasCount int        `json:"acknowledgedReplicasCount"`
}

func (r *Scheduler) IsPending() bool {
	return r.WaitTill != nil && r.WaitTill.UnixNano() >= time.Now().UnixNano()
}

func (r *Scheduler) IsAcked() bool {
	return r.AcknowledgedReplicasCount == r.ReplicasToWait
}

func (r *Scheduler) AwaitReplica(amount, msec int) {
	fmt.Println(amount, msec)
	r.WaitTimeoutMS = &msec
	r.ReplicasToWait = amount
	r.AcknowledgedReplicasCount = 0
	// set the suspend Time
	r.setDate(msec)
}

func (r *Scheduler) Release() {
	if r.AcknowledgedReplicasCount < r.ReplicasToWait {
		r.AcknowledgedReplicasCount += 1
	}
}

func (r *Scheduler) ToString() string {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	return string(b)
}

func (r *Scheduler) IncreasTotalReplicasCounter() {
	r.TotalReplicasCount++
}

func (r *Scheduler) setDate(msec int) {
	now := time.Now()
	waitTil := now.Add(time.Duration(msec) * time.Millisecond)
	r.WaitTill = &waitTil
}
