package core

import (
	"reflect"
	"testing"
)

func TestGetArgsMap(t *testing.T) {
	args := []string{"./your_program.sh", "--port", "9000", "--replicaof", "localhost 9000"}
	res := getArgumentMap(args)

	exp := map[Argument][]string{
		"--port":      []string{"9000"},
		"--replicaof": []string{"localhost", "9000"},
	}

	if !reflect.DeepEqual(exp, res) {
		t.Fail()
	}

	t.Log("OK")
}

func TestGetArgsMap2(t *testing.T) {
	args := []string{
		"./your_program.sh",
		"--port",
	}

	res := getArgumentMap(args)

	if len(res) != 0 {
		t.Fail()
	}

	t.Log("OK")
}