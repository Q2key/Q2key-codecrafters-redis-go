package core

import (
	"fmt"
	"github.com/codecrafters-io/redis-starter-go/app/contracts"
	"log"
	"os"
	"time"
)

type Instance struct {
	ReplicaId string
	Config    contracts.Config
	store     contracts.Store
}

const FakeReplicaId = "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb"

func NewRedisInstance(config contracts.Config) *Instance {
	ins := &Instance{
		ReplicaId: FakeReplicaId,
		store:     contracts.Store{},
		Config:    config,
	}

	ins.TryConnectDb()

	return ins
}

func (r *Instance) TryConnectDb() {
	if r.Config.GetDbFileName() == "" || r.Config.GetDir() == "" {
		return
	}

	path := fmt.Sprintf("%s/%s", r.Config.GetDir(), r.Config.GetDbFileName())

	db := NewRedisDB(path)
	if !db.IsFileExists(r.Config.GetDbFileName()) {
		_ = os.Mkdir(r.Config.GetDir(), os.ModeDir)
	}

	if !db.IsFileExists(path) {
		err := db.Create()
		if err != nil {
			log.Fatal(err)
		}
	}

	err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range db.GetData() {
		r.Set(k, v)
		exp, ok := db.GetExpires()[k]
		if ok {
			r.SetExpiredAt(k, exp)
		}
	}
}

func (r *Instance) Get(key string) contracts.Value {
	return r.store[key]
}

func (r *Instance) GetReplicaId() string {
	return r.ReplicaId
}

func (r *Instance) Set(key string, value string) {
	r.store[key] = &InstanceValue{
		Value: value,
	}
}

func (r *Instance) GetKeys(token string) []string {
	res := make([]string, 0)
	switch token {
	case "*":
		for k := range r.store {
			res = append(res, k)
		}
	}
	return res
}

func (r *Instance) GetStore() *map[string]contracts.Value {
	return &r.store
}

func (r *Instance) SetExpiredAt(key string, expired uint64) {
	tm := GetDateFromTimeStamp(expired)
	val, ok := r.store[key]
	if ok {
		val.SetExpired(tm)
	}

	r.store[key] = val
}

func (r *Instance) SetExpiredIn(key string, expiredIn uint64) {
	exp := time.Now().UTC().Add(time.Duration(expiredIn) * time.Millisecond)
	val, ok := r.store[key]
	if ok {
		val.SetExpired(exp)
	}
	r.store[key] = val
}

func (r *Instance) GetConfig() contracts.Config {
	return r.Config
}
