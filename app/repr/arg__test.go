package repr

import (
	"testing"
)

func TestToArgsShouldBeOk1(t *testing.T) {
	str := "*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n"
	_, res := ToArgs(str)
	if len(res) != 2 || res[0] != "ECHO" || res[1] != "hey" {
		t.Error()
	}

	t.Log("OK")
}

func TestToArgsShouldBeOk5(t *testing.T) {
	str := "*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n"
	_, res := ToArgs(str)
	if len(res) != 2 || res[0] != "ECHO" || res[1] != "hey" {
		t.Error()
	}

	t.Log("OK")
}

func TestToArgsShouldBeOk3(t *testing.T) {
	str := "*2\r\n$3\r\nGET\r\n$10\r\nstrawberry\r\n"
	_, res := ToArgs(str)
	if len(res) != 2 || res[0] != "GET" || res[1] != "strawberry" {
		t.Error()
	}

	t.Log("OK")
}

func TestToArgsShouldBeOk2(t *testing.T) {
	str := "*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n"
	_, res := ToArgs(str)
	if len(res) != 2 || res[0] != "ECHO" || res[1] != "hey" {
		t.Error()
	}

	t.Log("OK")
}
