package repr

import (
	"errors"
	"fmt"
	"strconv"
)

type ReprToken string

const (
	StringToken ReprToken = "$"
	ArrayToken  ReprToken = "*"
)

func StringsToRedisStrings(input []string) string {
	n := len(input)
	r := fmt.Sprintf("*%d", n)
	for i := 0; i < n; i++ {
		s := input[i]
		r += fmt.Sprintf("\r\n$%d\r\n%s", len(s), s)
	}
	r += "\r\n"
	return r
}

func FromStringToRedisBulkString(input string) string {
	l := len(input)
	return fmt.Sprintf("$%d\r\n%s\r\n", l, input)
}

func FromStringToRedisCommonString(input string) string {
	return fmt.Sprintf("+%s\r\n", input)
}

func FromStringToRedisInteger(str string) string {
	return fmt.Sprintf(":%s\r\n", str)
}

func ToRedisErrorString() string {
	return "$-1\r\n"
}

func GetValueType(q string) string {
	if len(q) == 0 {
		return "-"
	}

	if q[0] != '*' {
		return "-"
	}

	if ReprToken(q[1]) == StringToken {
		return "string"
	}

	if ReprToken(q[2]) == ArrayToken {
		return "array"
	}

	return "-"
}

func FromRedisStringToStringArray(q string) (error, []string) {
	if len(q) == 0 {
		return errors.New("empty string"), []string{}
	}

	// todo seems we have an error here :(
	if q[0] != '*' {
		return errors.New("invalid argument"), []string{}
	}

	sq := q[1:]
	n := len(sq)
	sli := make([]string, 0)
	for i := 0; i < n; i++ {
		if ReprToken(sq[i]) == StringToken {
			j := i + 1
			k := j

			for {
				ch := string(sq[k])
				if ch == "\r" {
					break
				} else {
					k += 1
				}
			}

			sl, err := strconv.Atoi(sq[j:k])
			if sl == 0 || err != nil {
				break
			}

			st := k + 2
			fi := st + sl

			if fi > len(sq) {
				break
			}

			sli = append(sli, sq[st:fi])
		}
	}

	return nil, sli
}
