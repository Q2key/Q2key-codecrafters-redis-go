package contracts

type Database interface {
	IsFileExists(string) bool
	Connect() error
	Create() error
	Data() map[string]string
	Expires() map[string]uint64
}
