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

type Argument string

const (
	Directory Argument = "--dir"
	Port      Argument = "--port"
	ReplicaOf Argument = "--replicaof"
)

func GetArgsMap(args []string) map[Argument][]string {
	n, m := len(args), map[Argument][]string{}
	for i := 0; i < n; i++ {
		j := i + 1
		if args[i][0] == '-' && j < n {
			m[Argument(args[i])] = strings.Split(args[j], " ")
		}
	}

	return m
}

func CreateConfigFromArgs(args []string) contracts.Config {
	c := core.NewConfig()
	m := GetArgsMap(args)

	val, ok := m[Directory]
	if ok && len(val) > 0 {
		c.SetDir(val[0])
	}

	val, ok = m[Port]
	if ok && len(val) > 0 {
		c.SetDir(val[0])
	}

	val, ok = m[ReplicaOf]
	if ok && len(val) == 2 {
		host, port := val[0], val[1]

		tcp := client.NewTcpClient(host, port)

		err := tcp.Connect()
		if err != nil {
			log.Fatal(err)
		}

		//Handshake 1
		tcp.SendBytes("*1\r\n$4\r\nPING\r\n")

		//Handshake 2
		req := FromStringsArray([]string{"REPLCONF", "listening-port", c.GetPort()})
		tcp.SendBytes(req)
		req = FromStringsArray([]string{"REPLCONF", "capa", "psync2"})
		tcp.SendBytes(req)

		//Handshake 3
		req = FromStringsArray([]string{"PSYNC", "?", "-1"})
		tcp.SendBytes(req)

		tcp.Disconnect()

		c.SetReplica(&contracts.Replica{
			OriginHost: host,
			OriginPort: port,
		})
	}

	return c
}

// todo need to rename | refactor
func ToArgs(q string) (error, []string) {
	if len(q) == 0 {
		return errors.New("empty string"), []string{}
	}

	//todo seems we have an error here :(
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
