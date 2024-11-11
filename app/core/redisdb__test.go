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
			t.Errorf(err.Error())
		}
	}(fp)

	db := NewRedisDB(fp)
	if db == nil {
		t.Error("db is nil")
		return
	}

	err := db.Create()
	if err != nil {
		t.Error(err)
		return
	}

	if !db.IsFileExists(fp) {
		t.Error("db is not exists")
		return
	}

	t.Log("OK!")
}

func TestConnectDatabaseShouldBeOk2(t *testing.T) {
	fp := "../../dump.rdb"
	db := NewRedisDB(fp)
	err := db.Connect()
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	if db.Data()["foo"] != "bar" {
		t.Error("foo is not bar")
	}

	if db.Data()["bas"] != "jazz" {
		t.Error("foo is not bar")
	}

	if db.Data()["nas"] != "lil" {
		t.Error("foo is not bar")
	}

	t.Log("OK!")

}
