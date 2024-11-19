package adapters

import "testing"

func TestFromStringShouldBeOk1(t *testing.T) {
	res := FromString("hello world")
	exp := "+hello world\r\n"
	if res != exp {
		t.Fail()
	}

	t.Log("OK!")
}

func TestFromStringShouldBeOk2(t *testing.T) {
	res := FromString("hello")
	exp := "+hello\r\n"
	if res != exp {
		t.Fail()
	}

	t.Log("OK!")
}

func TestFromStringArrayShouldBeOk1(t *testing.T) {
	res := FromStringsArray([]string{"hello", "world"})
	exp := "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"
	if res != exp {
		t.Fail()
	}

	t.Log("OK!")
}

func TestFromStringArrayShouldBeOk2(t *testing.T) {
	res := FromStringsArray([]string{"hello", "world", "my", "friend"})
	exp := "*4\r\n$5\r\nhello\r\n$5\r\nworld\r\n$2\r\nmy\r\n$6\r\nfriend\r\n"
	if res != exp {
		t.Fail()
	}

	t.Log("OK!")
}

func TestFromStringArrayShouldBeOk3(t *testing.T) {
	res := FromStringsArray([]string{""})
	exp := "*1\r\n$0\r\n\r\n"
	if res != exp {
		t.Fail()
	}

	t.Log("OK!")
}

func TestFromStringArrayShouldBeOk4(t *testing.T) {
	res := FromStringsArray([]string{})
	exp := "*0\r\n"
	if res != exp {
		t.Fail()
	}

	t.Log("OK!")
}

func TestBulkStringShouldBeOk1(t *testing.T) {
	res := BulkString("role:master")
	exp := "$11\r\nrole:master\r\n"
	if res != exp {
		t.Fail()
	}

	t.Log("OK!")
}

func TestBulkStringShouldBeOk2(t *testing.T) {
	res := FromString("FULLRESYNC 8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb 0")
	exp := "+FULLRESYNC 8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb 0\r\n"
	if res != exp {
		t.Fail()
	}

	t.Log("OK!")
}
