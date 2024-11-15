package config

import "github.com/codecrafters-io/redis-starter-go/app/contracts"

type Config struct {
	dir        string
	dbfilename string
	port       string
	replica    *contracts.Replica
}

const DefaultPort = "6379"

func NewConfig(dir, dbfilename string) *Config {
	return &Config{
		dir:        dir,
		dbfilename: dbfilename,
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
