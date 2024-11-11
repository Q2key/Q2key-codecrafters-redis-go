package contracts

type DataMap = map[string]string

type ExpiresMap = map[string]uint64

type Database interface {
	IsFileExists(string) bool
	Connect() error
	Create() error
	GetData() DataMap
	GetExpires() ExpiresMap
}
