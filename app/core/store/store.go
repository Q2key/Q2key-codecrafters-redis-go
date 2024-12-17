package store

import (
	"github.com/codecrafters-io/redis-starter-go/app/core/binary"
	"time"
)

type Store struct {
	kvs map[string]StoreValue
}

func NewStore() *Store {
	return &Store{kvs: make(map[string]StoreValue)}
}

func (r *Store) Get(key string) (StoreValue, bool) {
	val, ok := r.kvs[key]
	return val, ok
}

func (r *Store) Set(key string, value string) {
	r.kvs[key] = StoreValue{
		Value: value,
	}
}

func (r *Store) GetKeys(key string) []string {
	res := make([]string, 0)
	switch key {
	case "*":
		for k := range r.kvs {
			res = append(res, k)
		}
	}
	return res
}

func (r *Store) SetExpiredAt(key string, expiredAt uint64) {
	tm := binary.GetDateFromTimeStamp(expiredAt)
	val, ok := r.kvs[key]
	if ok {
		val.SetExpired(tm)
		r.kvs[key] = val
	}
}

func (r *Store) SetExpiredIn(key string, expiredIn uint64) {
	exp := time.Now().UTC().Add(time.Duration(expiredIn) * time.Millisecond)
	val, ok := r.kvs[key]
	if ok {
		val.SetExpired(exp)
		r.kvs[key] = val
	}
}
