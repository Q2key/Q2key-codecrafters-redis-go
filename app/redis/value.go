package redis

import "time"

type Value struct {
	Value   string
	Expired time.Time
	Created time.Time
}

func (rv *Value) IsExpired() bool {
	if rv.Expired.IsZero() {
		return false
	}

	now := time.Now()

	diff := rv.Expired.UnixMilli() - now.UnixMilli()

	return diff <= 0
}
