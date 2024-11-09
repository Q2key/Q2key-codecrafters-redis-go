package core

import (
	"testing"
)

func TestSetExpiring0(t *testing.T) {
	r := Instance{
		store: map[string]Value{},
	}

	r.Set("key", "value")

	v := r.Get("key")

	if !v.Expired.IsZero() {
		t.Fail()
	}

	t.Log("OK!")
}
