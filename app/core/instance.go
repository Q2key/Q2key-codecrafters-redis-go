package core

import (
	"fmt"
	"github.com/codecrafters-io/redis-starter-go/app/contracts"
	"github.com/codecrafters-io/redis-starter-go/app/rbyte"
	"log"
	"os"
	"time"
)

type Instance struct {
	Config contracts.Config
	store  contracts.Store
}

func NewRedisInstance(config contracts.Config) *Instance {
	ins := &Instance{
		store:  contracts.Store{},
		Config: config,
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
	tm := rbyte.GetDateFromTimeStamp(expired)
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
