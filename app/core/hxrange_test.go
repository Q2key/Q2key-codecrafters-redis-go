package core

import (
	"reflect"
	"testing"
)

func TestGetArgMap(t *testing.T) {
	if !reflect.DeepEqual(
		[][]string{{"mango", "0-0"}, {"pear", "0-1"}},
		getArgMap([]string{"xread", "streams", "mango", "pear", "0-0", "0-1"})) {
		t.FailNow()
	}

	t.Log("OK!")
}

func TestParseId(t *testing.T) {
	var ts float64
	var seq int

	ts, seq = parseId("0-0")
	if ts != 0 && seq != 0 {
		t.FailNow()
	}

	ts, seq = parseId("1-0")
	if ts != 1 && seq != 0 {
		t.FailNow()
	}

	ts, seq = parseId("1-1")
	if ts != 1 && seq != 1 {
		t.FailNow()
	}

	ts, seq = parseId("*")
	if ts != -1 && seq != -1 {
		t.FailNow()
	}

	t.Log("OK!")
}

func TestTryGetBlock(t *testing.T) {
	ok, _ := tryGetBlockParams([]string{"xread", "block", "1000", "streams", "mango", "0-1"})
	if !ok {
		t.Fail()
	}

	t.Log("OK!")
}
