package adapters

import (
	"errors"
	"github.com/codecrafters-io/redis-starter-go/app/client"
	"github.com/codecrafters-io/redis-starter-go/app/contracts"
	"github.com/codecrafters-io/redis-starter-go/app/core"
	"log"
	"strconv"
	"strings"
)

type ReprToken string

const (
	StringToken ReprToken = "$"
	ArrayToken  ReprToken = "*"
)

func GetArgsMap(args []string) map[string]string {
	argmap := make(map[string]string)

	return argmap
}

func ConfigFromArgs(args []string) contracts.Config {
	cfg := core.NewConfig("", "")

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

			if a == "--replicaof" && len(a) > 3 {
				parts := strings.Split(args[i+1], " ")
				replica := &contracts.Replica{
					OriginHost: parts[0],
					OriginPort: parts[1],
				}

				tcp := client.NewTcpClient(replica.OriginHost, replica.OriginPort)
				err := tcp.Connect()
				if err != nil {
					log.Fatal(err)
				}

				//step1
				_, _ = tcp.SendBytes("*1\r\n$4\r\nPING\r\n")

				//step2
				req := FromStringsArray([]string{"REPLCONF", "listening-port", cfg.GetPort()})
				_, _ = tcp.SendBytes(req)

				req = FromStringsArray([]string{"REPLCONF", "capa", "psync2"})
				_, _ = tcp.SendBytes(req)

				//step3
				req = FromStringsArray([]string{"PSYNC", "?", "-1"})
				_, _ = tcp.SendBytes(req)

				tcp.Disconnect()

				cfg.SetReplica(replica)
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
