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
		Config: &Config{
			dir:        "",
			dbfilename: "",
		},
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
		_ = os.Mkdir(ri.Config.dir, os.ModeDir)
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
	r.store[key] = Value{
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

func (r *Instance) GetStore() *map[string]Value {
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
