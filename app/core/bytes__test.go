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

	t.Log("OK!")
}

func TestParseMSecDateTimeStamp3(t *testing.T) {
	buff := []byte{0x00, 0x0c, 0x28, 0x8a, 0xc7, 0x01, 0x00, 0x00}
	ext := ParseMSecDateTimeStamp(&buff)
	if GetDateFromTimeStamp(ext).UTC().Format("2006-01-02 15:04:05") != "2032-01-01 00:00:00" {
		t.Fail()
	}

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

	ok, key, val := ParseValuePair(1, &buff)
	if !ok {
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
