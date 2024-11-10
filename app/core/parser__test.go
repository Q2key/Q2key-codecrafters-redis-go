package core

import (
	"fmt"
	"testing"
	"time"
)

func TestParseMSecDateTimeStamp(t *testing.T) {

	buff := []byte{0x15, 0x72, 0xE7, 0x07, 0x8F, 0x01, 0x00, 0x00}
	ext := ParseMSecDateTimeStamp(&buff)
	if ext != 1713824559637 {
		t.Fail()
	}

	tm := time.UnixMilli(int64(ext)).UTC()
	fmt.Println(tm)

	t.Log("OK!")
}

func TestParseMSecDateTimeStamp3(t *testing.T) {
	buff := []byte{0x00, 0x0c, 0x28, 0x8a, 0xc7, 0x01, 0x00, 0x00}
	ext := ParseMSecDateTimeStamp(&buff)

	tm1 := time.Unix(int64(ext)/1000, 0).UTC()
	tm2 := GetDateFromTimeStamp(ext)
	fmt.Println(tm1, tm2)

	t.Log("OK!")
}

func TestParseMSecDateTimeStamp2(t *testing.T) {

	tm := time.Unix(1956528000000/1000, 0).UTC()
	fmt.Println(tm)

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
