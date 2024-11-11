package repr

import (
	"strconv"
)

type ReprToken string

const (
	StringToken ReprToken = "$"
	ArrayToken  ReprToken = "*"
)

func ToArgs(q string) []string {
	s := q[1:]
	sli := make([]string, 0)
	for i := 0; i < len(s); i++ {
		if ReprToken(s[i]) == StringToken {
			j := i + 1
			k := j

			for {
				ch := string(s[k])
				if ch == "\r" {
					break
				} else {
					k += 1
				}
			}

			sl, err := strconv.Atoi(s[j:k])
			if sl == 0 || err != nil {
				break
			}

			st := k + 2
			fi := st + sl

			if fi > len(s) {
				break
			}

			sli = append(sli, s[st:fi])
		}
	}

	return sli
}
