package core

type Database interface {
	ReadFrom() (error, *Instance)
	Save(store *Instance) (error, *Instance)
	IsFileExists(name string) bool
	Connect() error
}
