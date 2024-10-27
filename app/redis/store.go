package redis

import (
	"time"
)

type Config struct {
	dir        string
	dbfilename string
}

func (r *Config) SetDir(val string) {
	r.dir = val
}

func (r *Config) SetDbFilenabe(val string) {
	r.dbfilename = val
}

func (r *Config) GetDir() string {
	return r.dir
}

func (r *Config) GetDbFileName() string {
	return r.dbfilename
}

type Store struct {
	store  map[string]Value
	config *Config
}

func NewRedisStore() Store {
	return Store{
		store: map[string]Value{},
	}
}

func (r *Store) GetConfig() *Config {
	return r.config
}

func (r *Store) SetConfig(config *Config) {
	r.config = config
}

func (r *Store) Get(key string) Value {
	return r.store[key]
}

func (r *Store) Set(key string, value string, expired int64) {
	now := time.Now().UTC()
	var exp time.Time
	if expired != 0 {
		exp = now.Add(time.Duration(expired) * time.Millisecond)
	}

	r.store[key] = Value{Value: value, Created: now, Expired: exp}
}
