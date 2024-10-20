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

func (rs *Store) Get(key string) Value {
	return rs.store[key]
}

func (rs *Store) Set(key string, value string, expired int64) {
	now := time.Now().UTC()
	var exp time.Time
	if expired != 0 {
		exp = now.Add(time.Duration(expired) * time.Millisecond)
	}

	rs.store[key] = Value{Value: value, Created: now, Expired: exp}
}

func (rs *Store) ToOkString(input string) []byte {
	return []byte(fmt.Sprintf("+%s\r\n", input))
}

func (rs *Store) ToErrorString() []byte {
	return []byte("$-1\r\n")
}
