package core

import (
	"fmt"
	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/contracts"
	"github.com/codecrafters-io/redis-starter-go/app/rbyte"
	"log"
	"os"
	"strings"
	"time"
)

type Instance struct {
	Config contracts.Config
	store  contracts.Store
}

func NewRedisInstance() *Instance {
	return &Instance{
		store:  contracts.Store{},
		Config: config.NewConfig("", ""),
	}
}

func NewRedisInstanceWithArgs(args []string) *Instance {
	ri := NewRedisInstance()
	if len(args) > 1 {
		for i := 1; i < len(args); i++ {
			a := args[i]
			if i+1 == len(args) {
				break
			}

			v := args[i+1]
			if a == "--dir" {
				ri.Config.SetDir(v)
			}

			if a == "--dbfilename" {
				ri.Config.SetDbFileName(v)
			}

			if a == "--port" {
				ri.Config.SetPort(v)
			}

			//todo think about validation
			if a == "--replicaof" && len(a) > 3 {
				v := args[i+1]
				parts := strings.Split(v, " ")
				fmt.Println("Replica host: ", parts[0])
				fmt.Println("Replica port: ", parts[1])
				ri.Config.SetReplica(&contracts.Replica{
					OriginHost: parts[0],
					OriginPort: parts[1],
				})
			}
		}
	}

	if ri.Config.GetDir() == "" {
		return ri
	}

	if ri.Config.GetDir() == "" {
		return ri
	}

	dbpath := fmt.Sprintf("%s/%s", ri.Config.GetDir(), ri.Config.GetDbFileName())

	db := NewRedisDB(dbpath)
	if !db.IsFileExists(ri.Config.GetDbFileName()) {
		_ = os.Mkdir(ri.Config.GetDir(), os.ModeDir)
	}

	if !db.IsFileExists(dbpath) {
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
		ri.Set(k, v)
		exp, ok := db.GetExpires()[k]
		if ok {
			ri.SetExpiredAt(k, exp)
		}
	}

	return ri
}

func (r *Instance) Get(key string) contracts.Value {
	return r.store[key]
}

func (r *Instance) Set(key string, value string) {
	r.store[key] = &Value{
		Value: value,
	}
}

func (r *Instance) Keys(token string) []string {
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
