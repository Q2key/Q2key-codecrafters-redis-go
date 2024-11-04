package core

import (
	"testing"
	"time"
)

func TestSetExpiring1000(t *testing.T) {
	r := Instance{
		store: map[string]Value{},
	}

	r.Set("key", "value", 1000)

	v := r.Get("key")

	t.Logf("Created: %s\r\n", v.Created.Format(time.RFC3339Nano))
	t.Logf("Expired: %s\r\n", v.Expired.Format(time.RFC3339Nano))

	d := v.Expired.UnixMilli() - v.Created.UnixMilli()
	if d != 1000 {
		t.Fail()
	}
}

func TestSetExpiring0(t *testing.T) {
	r := Instance{
		store: map[string]Value{},
	}

	r.Set("key", "value", 0)

	v := r.Get("key")
	t.Logf("Created: %s\r\n", v.Created.Format(time.RFC3339Nano))
	t.Logf("Expired: %s\r\n", v.Expired.Format(time.RFC3339Nano))

	if !v.Expired.IsZero() {
		t.Fail()
	}

	t.Log("OK!")
}
