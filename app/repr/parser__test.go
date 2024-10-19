package repr

import (
	"testing"
)

func TestParseArray(t *testing.T) {
	str := `*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n`
	res := ParseArray(str[1:])
	if len(res) != 2 && res[0] != "ECHO" && res[1] != "hey" {
		t.Error()
	}

}
