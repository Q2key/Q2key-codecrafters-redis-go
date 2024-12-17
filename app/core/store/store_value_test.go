package store

import (
	"testing"
	"time"
)

func TestGetSetValue(t *testing.T) {
	s := *NewStore()

	s.Set("Key0", "Value0")
	s.Set("Key1", "Value1")
	s.Set("Key2", "Value2")

	v1, _ := s.Get("Key0")
	if v1.GetValue() != "Value0" {
		t.Fail()
	}

	v2, _ := s.Get("Key1")
	if v2.GetValue() != "Value1" {
		t.Fail()
	}

	v3, _ := s.Get("Key2")
	if v3.GetValue() != "Value2" {
		t.Fail()
	}

	t.Log("OK!")
}

func TestShouldBeExpired2000(t *testing.T) {
	r := NewStore()

	r.Set("key", "value")
	r.SetExpiredIn("key", 2000)

	time.Sleep(3 * time.Second)

	v, _ := r.Get("key")
	if !v.IsExpired() {
		t.Fail()
	}

	t.Log("OK")
}

func TestShouldBeExpired100(t *testing.T) {
	r := NewStore()

	r.Set("key", "value")
	r.SetExpiredIn("key", 100)

	time.Sleep(101 * time.Millisecond)

	v, _ := r.Get("key")
	if !v.IsExpired() {
		t.Fail()
	}

	t.Log("OK")
}

func TestShouldBeExpired101(t *testing.T) {
	r := NewStore()

	r.Set("key", "value")
	r.SetExpiredIn("key", 101)

	time.Sleep(101 * time.Millisecond)

	v, _ := r.Get("key")
	if !v.IsExpired() {
		t.Fail()
	}

	t.Log("OK")
}

func TestShouldNotBeExpired0(t *testing.T) {
	r := NewStore()

	r.Set("key", "value")
	v, _ := r.Get("key")

	if v.IsExpired() {
		t.Fail()
	}

	if v.GetValue() != "value" {
		t.Fail()
	}

	t.Log("OK!")
}

func TestShouldNotBeExpired4000(t *testing.T) {
	r := NewStore()
	r.Set("key", "value")
	r.SetExpiredIn("key", 4000)
	v, _ := r.Get("key")

	time.Sleep(3 * time.Second)

	if v.IsExpired() {
		t.Fail()
	}

	if v.GetValue() != "value" {
		t.Fail()
	}

	t.Log("OK!")
}
