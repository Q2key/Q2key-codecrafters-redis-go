package contracts

type Store interface {
	Get(string) (Value, bool)
	Set(string, string)
	GetKeys(string) []string
	SetExpiredAt(key string, expiredAt uint64)
	SetExpiredIn(key string, expiredIn uint64)
}
