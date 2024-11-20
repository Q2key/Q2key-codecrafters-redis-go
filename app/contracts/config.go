package contracts

type Config interface {
	SetDir(string)
	SetDbFileName(string)
	SetReplica(replica *Replica)
	SetPort(string)
	GetDir() string
	GetDbFileName() string
	GetPort() string
	GetReplica() *Replica
	FromArguments([]string) *Config
}
