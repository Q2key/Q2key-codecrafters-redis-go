package main

import (
	"testing"
	"time"
)

func TestShouldBeExpired2000(t *testing.T) {

	r := RedisInstance{
		Store: map[string]RedisValue{},
	}

	r.SetExpiring("key", "value", 2000)

	time.Sleep(3 * time.Second)

	v := r.Get("key")
	if !r.IsExpired(v) {
		t.Fail()
	}

	t.Log("OK")
}

func TestShouldBeExpired100(t *testing.T) {

	r := RedisInstance{
		Store: map[string]RedisValue{},
	}

	r.SetExpiring("key", "value", 100)

	time.Sleep(101 * time.Millisecond)

	v := r.Get("key")
	if !r.IsExpired(v) {
		t.Fail()
	}

	t.Log("OK")
}

func TestShouldNotBeExpired1(t *testing.T) {

	r := RedisInstance{
		Store: map[string]RedisValue{},
	}

	r.SetExpiring("key", "value", 0)
	v := r.Get("key")

	if r.IsExpired(v) {
		t.Fail()
		return
	}

	if v.Value != "value" {
		t.Fail()

	}

	t.Log("OK!")
}

func TestShouldNotBeExpired2(t *testing.T) {

	r := RedisInstance{
		Store: map[string]RedisValue{},
	}

	r.SetExpiring("key", "value", 4000)
	v := r.Get("key")

	time.Sleep(3 * time.Second)

	if r.IsExpired(v) {
		t.Fail()

	}

	if v.Value != "value" {
		t.Fail()
	}

	t.Log("OK!")
}

func TestSetExpiring1000(t *testing.T) {

	r := RedisInstance{
		Store: map[string]RedisValue{},
	}

	r.SetExpiring("key", "value", 1000)

	v := r.Get("key")

	t.Logf("Created: %s\r\n", v.Created.Format(time.RFC3339Nano))
	t.Logf("Expired: %s\r\n", v.Expired.Format(time.RFC3339Nano))

	d := v.Expired.UnixMilli() - v.Created.UnixMilli()
	if d != 1000 {
		t.Fail()
	}
}

func TestSetExpiring0(t *testing.T) {

	r := RedisInstance{
		Store: map[string]RedisValue{},
	}

	r.SetExpiring("key", "value", 0)

	v := r.Get("key")
	t.Logf("Created: %s\r\n", v.Created.Format(time.RFC3339Nano))
	t.Logf("Expired: %s\r\n", v.Expired.Format(time.RFC3339Nano))

	if !v.Expired.IsZero() {
		t.Fail()
	}

	t.Log("OK!")
}
