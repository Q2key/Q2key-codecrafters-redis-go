package contracts

type DBFileConnector interface {
	IsFileExists(string) bool
	Connect() error
	Create() error
	GetData() map[string]string
	GetExpires() map[string]uint64
}
