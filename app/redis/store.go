package redis

import (
	"fmt"
	"time"
)

type Store struct {
	store map[string]Value
}

func NewRedisStore() Store {
	return Store{store: map[string]Value{}}
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

func (r *Store) ToOkString(input string) []byte {
	return []byte(fmt.Sprintf("+%s\r\n", input))
}

func (r *Store) ToErrorString() []byte {
	return []byte("$-1\r\n")
}
