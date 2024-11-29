package core

import (
	"testing"

	"github.com/codecrafters-io/redis-starter-go/app/contracts"
)

func TestSetExpiring0(t *testing.T) {
	r := Instance{
		Store: map[string]contracts.Value{},
	}

	r.Set("key", "value")

	v := r.Get("key")

	if v.IsExpired() {
		t.Error("key should not be expired")
	}

	t.Log("OK!")
}
