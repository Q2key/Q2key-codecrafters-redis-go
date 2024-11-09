package core

import (
	"testing"
	"time"
)

func TestGetSetValue(t *testing.T) {
	s := NewRedisInstance()

	s.Set("Key0", "Value0")
	s.Set("Key1", "Value1")
	s.Set("Key2", "Value2")

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

	r := Instance{
		store: map[string]Value{},
	}

	r.Set("key", "value")
	r.SetExpiredIn("key", 2000)

	time.Sleep(3 * time.Second)

	v := r.Get("key")
	if !v.IsExpired() {
		t.Fail()
	}

	t.Log("OK")
}

func TestShouldBeExpired100(t *testing.T) {

	r := Instance{
		store: map[string]Value{},
	}

	r.Set("key", "value")
	r.SetExpiredIn("key", 100)

	time.Sleep(101 * time.Millisecond)

	v := r.Get("key")
	if !v.IsExpired() {
		t.Fail()
	}

	t.Log("OK")
}

func TestShouldBeExpired101(t *testing.T) {

	r := Instance{
		store: map[string]Value{},
	}

	r.Set("key", "value")
	r.SetExpiredIn("key", 101)

	time.Sleep(101 * time.Millisecond)

	v := r.Get("key")
	if !v.IsExpired() {
		t.Fail()
	}

	t.Log("OK")
}

func TestShouldNotBeExpired0(t *testing.T) {

	r := Instance{
		store: map[string]Value{},
	}

	r.Set("key", "value")
	v := r.Get("key")

	if v.IsExpired() {
		t.Fail()
	}

	if v.Value != "value" {
		t.Fail()

	}

	t.Log("OK!")
}

func TestShouldNotBeExpired4000(t *testing.T) {

	r := Instance{
		store: map[string]Value{},
	}

	r.Set("key", "value")
	r.SetExpiredIn("key", 4000)
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

func TestReturnCorrectConfig(t *testing.T) {
	s := &Instance{
		store: make(map[string]Value),
		Config: &Config{
			dir:        "temp",
			dbfilename: "develop",
		},
	}

	s.Config.SetDbFileName("develop")
	if s.Config.dbfilename != "develop" {
		t.Fail()
	}

	if s.Config.dir != "temp" {
		t.Fail()
	}

	t.Log("OK!")
}
