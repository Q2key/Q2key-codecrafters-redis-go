package repr

import (
	"testing"
)

func TestTestLen(t *testing.T) {
	n := TestLen("\rHELLO")
	if n != 6 {
		t.Error()
	}

	n = TestLen("*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n")
	if n != 23 {
		t.Error()
	}

	t.Log("OK")
}

func TestParseArray(t *testing.T) {
	str := "*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n"
	res := ParseArray(str[1:])
	if len(res) != 2 || res[0] != "ECHO" || res[1] != "hey" {
		t.Error()
	}

	t.Log("OK")
}

func TestToStrigArray(t *testing.T) {
	inp := []string{"dir", "temp"}
	res := string(ToStringArray(inp))
	// *<number-of-elements>\r\n<element-1>...<element-n>

	exp := "*2\r\n$3\r\ndir\r\n$4\r\ntemp\r\n"
	if res != exp {
		t.Errorf("Got %s", res)
		t.Fail()
	}

	t.Log("OK")
}
