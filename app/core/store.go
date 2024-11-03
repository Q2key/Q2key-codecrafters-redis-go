package core

import (
	"log"
	"time"
)

type Instance struct {
	Config *Config
	store  map[string]Value
}

func NewRedisInstance() *Instance {
	return &Instance{
		store: map[string]Value{},
	}
}

func (r *Instance) WithArgs(args []string) *Instance {
	if r.store == nil {
		log.Fatal("No Instance created")
	}

	if len(args) > 1 {
		for i := 1; i < len(args); i++ {
			a := args[i]
			if i+1 == len(args) {
				break
			}

			v := args[i+1]
			if a == "--dir" {
				r.Config.SetDir(v)
			}

			if a == "--dbfilename" {
				r.Config.SetDbFileName(v)
			}
		}
	}

	return r
}

func (r *Instance) Get(key string) Value {
	return r.store[key]
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

func (r *Instance) Set(key string, value string, expired int64) {
	now := time.Now().UTC()
	var exp time.Time
	if expired != 0 {
		exp = now.Add(time.Duration(expired) * time.Millisecond)
	}

	r.store[key] = Value{Value: value, Created: now, Expired: exp}
}
