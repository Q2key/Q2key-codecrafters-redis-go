package redis

import "testing"

func TestReturnCorrectConfig(t *testing.T) {
	s := &Store{
		store: make(map[string]Value),
		config: &Config{
			dir:        "temp",
			dbfilename: "develop",
		},
	}

	cfg := s.GetConfig()
	if cfg.dbfilename != "develop" {
		t.Fail()
	}

	if cfg.dir != "temp" {
		t.Fail()
	}

	t.Log("OK!")
}

func TestSetCorrectConfig(t *testing.T) {
	s := &Store{
		store: make(map[string]Value),
		config: &Config{
			dir:        "temp",
			dbfilename: "develop",
		},
	}

	s.SetConfig(&Config{
		dir:        "/var/log",
		dbfilename: "production",
	})

	c := s.GetConfig()
	if c.dir != "/var/log" {
		t.Fail()
	}

	if c.dbfilename != "production" {
		t.Fail()
	}

	t.Log("OK")
}
