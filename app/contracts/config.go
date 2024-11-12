package contracts

type Config interface {
	SetDir(string)
	SetDbFileName(string)
	SetPort(string)
	GetDir() string
	GetDbFileName() string
	GetPort() string
}
