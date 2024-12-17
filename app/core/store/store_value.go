package store

import (
	"time"
)

type StoreValue struct {
	Value   string
	Expired *time.Time
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