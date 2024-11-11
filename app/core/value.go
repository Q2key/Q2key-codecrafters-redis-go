package core

import (
	"time"
)

type Value struct {
	Value   string
	Expired *time.Time
}

func (r *Value) IsExpired() bool {
	if r.Expired == nil {
		return false
	}

	return r.Expired.UnixNano() <= time.Now().UTC().UnixNano()
}

func (r *Value) SetExpired(expired time.Time) {
	r.Expired = &expired
}

func (r *Value) GetValue() string {
	return r.Value
}
