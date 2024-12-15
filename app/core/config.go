package core

import (
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/contracts"
)

type Config struct {
	dir        string
	dbfilename string
	port       string
	replica    *contracts.Replica
	isMaster   bool
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
	Directory  Argument = "--dir"
	DbFilename Argument = "--dbfilename"
	Port       Argument = "--port"
	ReplicaOf  Argument = "--replicaof"
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

func createConfigFromArgs(args []string) contracts.Config {
	c := NewConfig()
	m := getArgumentMap(args)

	val, ok := m[Directory]
	if ok && len(val) > 0 {
		c.SetDir(val[0])
	}

	val, ok = m[Port]
	if ok && len(val) > 0 {
		c.SetPort(val[0])
	}

	val, ok = m[DbFilename]
	if ok && len(val) > 0 {
		c.SetDbFileName(val[0])
	}

	val, ok = m[ReplicaOf]
	if ok && len(val) >= 2 {
		c.SetReplica(contracts.NewReplica(val[0], val[1]))
	}

	return c
}

func (r *Config) IsMaster() bool {
	return r.replica == nil
}
