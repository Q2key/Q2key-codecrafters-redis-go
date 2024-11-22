package core

import (
	"github.com/codecrafters-io/redis-starter-go/app/client"
	"github.com/codecrafters-io/redis-starter-go/app/contracts"
	"log"
	"strings"
)

type Config struct {
	dir        string
	dbfilename string
	port       string
	replica    *contracts.Replica
}

const DefaultPort = "6379"

func NewConfig() *Config {
	return &Config{
		dir:        "",
		dbfilename: "",
		port:       DefaultPort,
	}
}

func (r *Config) SetDir(val string) {
	r.dir = val
}

func (r *Config) SetPort(val string) {
	r.port = val
}

func (r *Config) SetReplica(val *contracts.Replica) {
	r.replica = val
}

func (r *Config) SetDbFileName(val string) {
	r.dbfilename = val
}

func (r *Config) GetDir() string { return r.dir }

func (r *Config) GetDbFileName() string { return r.dbfilename }

func (r *Config) GetPort() string { return r.port }

func (r *Config) GetReplica() *contracts.Replica { return r.replica }

type Argument string

const (
	Directory Argument = "--dir"
	Port      Argument = "--port"
	ReplicaOf Argument = "--replicaof"
)

// FromArguments todo change dep
func (r *Config) FromArguments(args []string) *contracts.Config {
	cfg := createConfigFromArgs(args)
	return &cfg
}

func getArgumentMap(args []string) map[Argument][]string {
	n, m := len(args), map[Argument][]string{}
	for i := 0; i < n; i++ {
		j := i + 1
		if args[i][0] == '-' && j < n {
			m[Argument(args[i])] = strings.Split(args[j], " ")
		}
	}

	return m
}

func initHandShake(c contracts.Config) {
	tcp := client.NewTcpClient(c.GetReplica().OriginPort, c.GetReplica().OriginHost)

	err := tcp.Connect()
	if err != nil {
		log.Fatal(err)
	}

	//Handshake 1
	tcp.SendBytes("*1\r\n$4\r\nPING\r\n")

	//Handshake 2
	req := FromStringArrayToRedisStringArray([]string{"REPLCONF", "listening-port", c.GetPort()})
	tcp.SendBytes(req)
	req = FromStringArrayToRedisStringArray([]string{"REPLCONF", "capa", "psync2"})
	tcp.SendBytes(req)

	//Handshake 3
	req = FromStringArrayToRedisStringArray([]string{"PSYNC", "?", "-1"})
	tcp.SendBytes(req)

	tcp.Disconnect()
}

func createConfigFromArgs(args []string) contracts.Config {
	c := NewConfig()
	m := getArgumentMap(args)

	val, ok := m[Directory]
	if ok && len(val) > 0 {
		c.SetDir(val[0])
	}

	val, ok = m[Port]
	if ok && len(val) > 0 {
		c.SetDir(val[0])
	}

	val, ok = m[ReplicaOf]
	if ok && len(val) >= 2 {
		c.SetReplica(contracts.NewReplica(val[0], val[1]))
	}

	if c.replica != nil {
		initHandShake(c)
	}

	return c
}
