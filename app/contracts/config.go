package contracts

type Config interface {
	SetDir(string)
	SetDbFileName(string)
	GetDir() string
	GetDbFileName() string
}
