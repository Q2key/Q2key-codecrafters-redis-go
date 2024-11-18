package repr

import (
	"errors"
	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/contracts"
	"strconv"
	"strings"
)

type ReprToken string

const (
	StringToken ReprToken = "$"
	ArrayToken  ReprToken = "*"
)

func ConfigFromArgs(args []string) contracts.Config {
	cfg := config.NewConfig("", "")

	if len(args) > 1 {
		for i := 1; i < len(args); i++ {
			a := args[i]
			if i+1 == len(args) {
				break
			}

			v := args[i+1]
			if a == "--dir" {
				cfg.SetDir(v)
			}

			if a == "--dbfilename" {
				cfg.SetDbFileName(v)
			}

			if a == "--port" {
				cfg.SetPort(v)
			}

			//todo think about validation
			if a == "--replicaof" && len(a) > 3 {
				v := args[i+1]
				parts := strings.Split(v, " ")
				cfg.SetReplica(&contracts.Replica{
					OriginHost: parts[0],
					OriginPort: parts[1],
				})
			}
		}
	}

	return cfg
}

func ToArgs(q string) (error, []string) {

	if len(q) == 0 {
		return errors.New("empty string"), []string{}
	}

	if q[0] != '*' {
		return errors.New("invalid argument"), []string{}
	}

	s := q[1:]
	n := len(s)
	sli := make([]string, 0)
	for i := 0; i < n; i++ {
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

	return nil, sli
}
