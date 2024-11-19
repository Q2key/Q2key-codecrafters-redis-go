package contracts

type Instance interface {
	Get(string) Value
	GetConfig() Config
	GetKeys(string) []string
	GetReplicaId() string
	Set(key string, value string)
	GetStore() *map[string]Value
	SetExpiredAt(string, uint64)
	SetExpiredIn(string, uint64)
}
