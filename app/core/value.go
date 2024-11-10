package core

import (
	"fmt"
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

	now := time.Now().UTC()

	isExpired := r.Expired.UnixNano() <= now.UnixNano()

	fmt.Printf("\t\r->Expired value: %s\n", r.Value)
	fmt.Printf("\t\r->Expired time: %v\n", r.Expired)
	fmt.Printf("\t\r->Now time: %s\n", now)
	fmt.Printf("\t\r->IsExpired: %v\n", isExpired)

	return isExpired

}

func (r *Value) SetExpired(expired time.Time) {
	r.Expired = &expired
}
