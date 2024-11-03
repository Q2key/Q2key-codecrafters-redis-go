package core

type Database interface {
	Name() string
	Path() string
	Create(name string) error
	ReadFrom() (error, *Instance)
	Save(store *Instance) (error, *Instance)
	IsFileExists(name string) bool
	Connect(name string) error
}
