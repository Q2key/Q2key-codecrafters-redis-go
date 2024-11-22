package contracts

type Replica struct {
	OriginHost string
	OriginPort string
}

func NewReplica(originHost string, originPort string) *Replica {
	return &Replica{originHost, originPort}
}
