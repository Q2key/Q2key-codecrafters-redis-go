package redis

import (
	"testing"
	"time"
)

func TestGetSetValue(t *testing.T) {
	s := NewRedisStore()

	s.Set("Key0", "Value0", 0)
	s.Set("Key1", "Value1", 0)
	s.Set("Key2", "Value2", 0)

	if s.Get("Key0").Value != "Value0" {
		t.Fail()
	}

	if s.Get("Key1").Value != "Value1" {
		t.Fail()
	}

	if s.Get("Key2").Value != "Value2" {
		t.Fail()
	}

	t.Log("OK!")
}

func TestShouldBeExpired2000(t *testing.T) {

	r := Store{
		store: map[string]Value{},
	}

	r.Set("key", "value", 2000)

	time.Sleep(3 * time.Second)

	v := r.Get("key")
	if !v.IsExpired() {
		t.Fail()
	}

	t.Log("OK")
}

func TestShouldBeExpired100(t *testing.T) {

	r := Store{
		store: map[string]Value{},
	}

	r.Set("key", "value", 100)

	time.Sleep(101 * time.Millisecond)

	v := r.Get("key")
	if !v.IsExpired() {
		t.Fail()
	}

	t.Log("OK")
}

func TestShouldBeExpired101(t *testing.T) {

	r := Store{
		store: map[string]Value{},
	}

	r.Set("key", "value", 101)

	time.Sleep(101 * time.Millisecond)

	v := r.Get("key")
	if !v.IsExpired() {
		t.Fail()
	}

	t.Log("OK")
}

func TestShouldNotBeExpired0(t *testing.T) {

	r := Store{
		store: map[string]Value{},
	}

	r.Set("key", "value", 0)
	v := r.Get("key")

	if v.IsExpired() {
		t.Fail()
		return
	}

	if v.Value != "value" {
		t.Fail()

	}

	t.Log("OK!")
}

func TestShouldNotBeExpired4000(t *testing.T) {

	r := Store{
		store: map[string]Value{},
	}

	r.Set("key", "value", 4000)
	v := r.Get("key")

	time.Sleep(3 * time.Second)

	if v.IsExpired() {
		t.Fail()

	}

	if v.Value != "value" {
		t.Fail()
	}

	t.Log("OK!")
}
