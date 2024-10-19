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

func ParseArray(q string) []string {
	vs := q[1:]
	sli := make([]string, 0)
	for i := 0; i < len(vs); i++ {
		if vs[i:i+1] == "$" {
			ni := i + 1
			nj := i + 2

			if nj > len(vs) {
				break
			}

			sl, err := strconv.Atoi(vs[ni:nj])
			if sl == 0 || err != nil {
				break
			}

			st := nj + 4
			fi := st + sl

			if fi > len(vs) {
				break
			}

			sli = append(sli, vs[st:fi])
		}
	}
	return sli
}
