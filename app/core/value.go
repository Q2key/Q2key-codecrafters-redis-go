package core

import "time"

type Value struct {
	Value   string
	Expired time.Time
	Created time.Time
}

func (r *Value) IsExpired() bool {
	if r.Expired.IsZero() {
		return false
	}

	now := time.Now()

	diff := r.Expired.UnixMilli() - now.UnixMilli()

	return diff <= 0
}
