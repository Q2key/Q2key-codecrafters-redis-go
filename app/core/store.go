package core

import (
	"time"
)

type Store struct {
	kvs map[string]*StoreValue
}

func NewStore() *Store {
	return &Store{kvs: make(map[string]*StoreValue)}
}

func (r *Store) Get(key string) (*StoreValue, bool) {
	val, ok := r.kvs[key]
	return val, ok
}

func (r *Store) Set(key string, value string) {
	r.kvs[key] = &StoreValue{
		Value:     value,
		ValueType: "string",
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
	tm := GetDateFromTimeStamp(expiredAt)
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

type StoreValue struct {
	Value     string
	Expired   *time.Time
	ValueType string
}

func (r *StoreValue) IsExpired() bool {
	if r.Expired == nil {
		return false
	}

	return r.Expired.UnixNano() <= time.Now().UTC().UnixNano()
}

func (r *StoreValue) SetExpired(expired time.Time) {
	r.Expired = &expired
}

func (r *StoreValue) GetValue() string {
	return r.Value
}
