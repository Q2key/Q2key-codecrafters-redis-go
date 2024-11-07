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

	fmt.Printf("\r\narguments: %v", args)
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

	//todo move to db
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
		exp, ok := db.Expires[k]
		if !ok {
			exp = 0
		}

		ri.Set(k, v, exp)
	}

	return ri
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

func (r *Instance) Store() *map[string]Value {
	return &r.store
}

func (r *Instance) Set(key string, value string, expired uint64) {
	now := time.Now().UTC()
	var exp time.Time
	if expired != 0 {
		exp = now.Add(time.Duration(expired) * time.Millisecond)
	}
	v := &Value{
		Val: value, Created: now, Expired: exp,
	}
	r.store[key] = *v
}
