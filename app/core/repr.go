package core

import (
	"fmt"
)

type REPRToken string

const (
	IntegerToken      REPRToken = ":"
	SimpleErrorToken  REPRToken = "-"
	SimpleStringToken REPRToken = "+"
	BulkStringToken   REPRToken = "$"
	ArrayToken        REPRToken = "*"
)

const (
	CLRF = "\r\n"
)

func ToRedisStrings(str []string) string {
	n := len(str)
	r := fmt.Sprintf("%s%d", ArrayToken, n)
	for i := 0; i < n; i++ {
		s := str[i]
		r += fmt.Sprintf("\r\n$%d\r\n%s", len(s), s)
	}
	r += "\r\n"
	return r
}

func ToRedisBulkString(str string) string {
	return fmt.Sprintf("%s%d\r\n%s\r\n", BulkStringToken, len(str), str)
}

func ToArrayDefString(n int) string {
	return fmt.Sprintf("%s%d\r\n", ArrayToken, n)
}

func ToRedisSimpleString(str string) string {
	return fmt.Sprintf("%s%s\r\n", SimpleStringToken, str)
}

func ToRedisInteger(str string) string {
	return fmt.Sprintf("%s%s\r\n", IntegerToken, str)
}

func ToRedisNullBulkString() string {
	return fmt.Sprintf("%s%s1\r\n", BulkStringToken, SimpleErrorToken)
}

func ToSimpleError(str string) string {
	return fmt.Sprintf("%s%s\r\n", SimpleErrorToken, str)
}
