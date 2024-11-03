package repr

import (
	"testing"
)

func TestToArgsShouldBeOk1(t *testing.T) {
	str := "*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n"
	res := ToArgs(str[1:])
	if len(res) != 2 || res[0] != "ECHO" || res[1] != "hey" {
		t.Error()
	}

	t.Log("OK")
}

func TestToArgsShouldBeOk2(t *testing.T) {
	str := "*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n"
	res := ToArgs(str[1:])
	if len(res) != 2 || res[0] != "ECHO" || res[1] != "hey" {
		t.Error()
	}

	t.Log("OK")
}
