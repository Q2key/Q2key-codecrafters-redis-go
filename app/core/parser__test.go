package core

import "testing"

func TestParseMSecDateTimeStamp(t *testing.T) {

	buff := []byte{0x15, 0x72, 0xE7, 0x07, 0x8F, 0x01, 0x00, 0x00}
	ext := ParseMSecDateTimeStamp(&buff)
	if ext != 1713824559637 {
		t.Fail()
	}

	t.Log("OK!")
}

func TestParseValuePair(t *testing.T) {

	buff := []byte{
		0x00,
		0x03, 0x62, 0x61, 0x7A,
		0x03, 0x71, 0x75, 0x78,
	}

	bool, key, val := ParseValuePair(1, &buff)
	if !bool {
		t.Fail()
	}

	if *key != "baz" {
		t.Fail()
	}

	if *val != "qux" {
		t.Fail()
	}

	t.Log("OK!")
}
