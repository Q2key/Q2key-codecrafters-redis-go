package core

import (
	"os"
	"testing"
)

func TestCreateDatabaseShouldBeOk1(t *testing.T) {
	fp := "./my.test.rdb"
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.FailNow()
		}
	}(fp)

	db := NewRedisDB(fp)
	if db == nil {
		t.FailNow()
		return
	}

	err := db.Create()
	if err != nil {
		t.FailNow()
		return
	}

	if !db.IsFileExists(fp) {
		t.FailNow()
	}

	t.Log("OK!")
}

func TestConnectDatabaseShouldBeOk2(t *testing.T) {
	fp := "../../dump.rdb"
	db := NewRedisDB(fp)
	err := db.Connect()
	if err != nil {
		t.FailNow()
	}

	if db.GetData()["foo"] != "bar" {
		t.FailNow()
	}

	if db.GetData()["bas"] != "jazz" {
		t.FailNow()
	}

	if db.GetData()["nas"] != "lil" {
		t.FailNow()
	}

	t.Log("OK!")
}
