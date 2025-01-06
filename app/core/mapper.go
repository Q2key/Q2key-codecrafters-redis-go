package core

import (
	"strconv"
)

func FromRedisStringToStringArray(q string) []string {
	if len(q) == 0 {
		return []string{}
	}

	// todo seems we have an error here :(
	if q[0] != '*' {
		return []string{}
	}

	sq := q[1:]
	n := len(sq)
	sli := make([]string, 0)
	for i := 0; i < n; i++ {
		if REPRToken(sq[i]) == BulkStringToken {
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

	return sli
}
