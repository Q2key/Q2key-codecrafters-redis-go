package repr

import (
	"strconv"
)

const (
	NA = iota
	Error
	String
	BulkStrings
	Number
	Array
)

func ParseType(q string) int {
	fb := q[0:1]
	//error
	switch fb {
	case "*":
		return Array
	case ":":
		return Number
	case "$":
		return BulkStrings
	default:
		return NA
	}
}

func TestLen(q string) int {
	return len(q)
}

func ParseArray(q string) []string {
	s := q[1:]
	sli := make([]string, 0)
	for i := 0; i < len(s); i++ {
		if string(s[i]) == "$" {
			ni := i + 1
			nj := i + 2

			sl, err := strconv.Atoi(s[ni:nj])
			if sl == 0 || err != nil {
				break
			}

			st := nj + 2
			fi := st + sl

			if fi > len(s) {
				break
			}

			sli = append(sli, s[st:fi])
		}
	}

	return sli
}

func ParseNumberValue(q string) string {
	return ""
}
