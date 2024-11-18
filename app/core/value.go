package core

import (
	"time"
)

type InstanceValue struct {
	Value   string
	Expired *time.Time
}

func (r *InstanceValue) IsExpired() bool {
	if r.Expired == nil {
		return false
	}

	return r.Expired.UnixNano() <= time.Now().UTC().UnixNano()
}

func (r *InstanceValue) SetExpired(expired time.Time) {
	r.Expired = &expired
}

func (r *InstanceValue) GetValue() string {
	return r.Value
}
