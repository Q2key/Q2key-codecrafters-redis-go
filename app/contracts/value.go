package contracts

import "time"

type Value interface {
	IsExpired() bool
	SetExpired(time time.Time)
	GetValue() string
}
