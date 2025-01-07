package core

import "time"

type RedisValue interface {
	ToString() string
	GetType() string
	SetValue(interface{})
	SetExpired(expired time.Time)
	IsExpired() bool
}
