package core

import (
	"fmt"
	"log"
	"os"
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

func NewRedisInstanceWithArgs(args []string) *Instance {
	ri := NewRedisInstance()
	ri.Config = &Config{
		dir:        "",
		dbfilename: "",
	}

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
		}
	}

	if ri.Config.dbfilename == "" {
		return ri
	}

	if ri.Config.dir == "" {
		return ri
	}

	dbpath := fmt.Sprintf("%s/%s", ri.Config.dir, ri.Config.dbfilename)

	db := NewRedisDB(dbpath)

	if !db.IsFileExists(ri.Config.dbfilename) {
		os.Mkdir(ri.Config.dir, os.ModeDir)
	}

	if !db.IsFileExists(dbpath) {
		db.Create()
	}

	err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range db.Data {
		ri.Set(k, v)
		exp, ok := db.Expires[k]
		if ok {
			ri.SetExpiredAt(k, exp)
		}
	}

	return ri
}

func (r *Instance) Get(key string) Value {
	return r.store[key]
}

func (r *Instance) Set(key string, value string) {
	var e time.Time
	r.store[key] = Value{
		Value: value, Expired: e,
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

func (r *Instance) Store() *map[string]Value {
	return &r.store
}

func (r *Instance) SetExpiredAt(key string, expired uint64) {
	tm := time.UnixMilli(int64(expired)).UTC()
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
