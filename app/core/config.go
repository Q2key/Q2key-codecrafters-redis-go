package core

import (
	"strings"
)

type ReplicationProps struct {
	OriginHost string
	OriginPort string
}

func NewReplicationProps(originHost string, originPort string) *ReplicationProps {
	return &ReplicationProps{originHost, originPort}
}

type Config struct {
	Dir        string
	DbFileName string
	Port       string
	Replica    *ReplicationProps
}

const DefaultPort = "6379"

func NewConfig() *Config {
	return &Config{
		Dir:        "",
		DbFileName: "",
		Port:       DefaultPort,
	}
}

func (r *Config) GetReplica() *ReplicationProps { return r.Replica }

type Argument string

const (
	Directory  Argument = "--dir"
	DbFilename Argument = "--dbfilename"
	Port       Argument = "--port"
	ReplicaOf  Argument = "--replicaof"
)

func (r *Config) FromArguments(args []string) *Config {
	return createConfigFromArgs(args)
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

func createConfigFromArgs(args []string) *Config {
	c := NewConfig()
	m := getArgumentMap(args)

	val, ok := m[Directory]
	if ok && len(val) > 0 {
		c.Dir = val[0]
	}

	val, ok = m[Port]
	if ok && len(val) > 0 {
		c.Port = val[0]
	}

	val, ok = m[DbFilename]
	if ok && len(val) > 0 {
		c.DbFileName = val[0]
	}

	val, ok = m[ReplicaOf]
	if ok && len(val) >= 2 {
		c.Replica = NewReplicationProps(val[0], val[1])
	}

	return c
}

func (r *Config) IsMaster() bool {
	return r.Replica == nil
}
