package contracts

type Ack struct {
	ConnId string
	Offset int
}

type Instance interface {
	InitHandshakeWithMaster()
	GetWrittenBytes() int
	GetConfig() Config
	GetMasterConn() RedisConn
	GetReplicas() map[string]RedisConn
	GetStore() Store
	GetAckChan() *chan Ack
	RegisterMasterConn(conn *RedisConn)
	RegisterReplicaConn(conn *RedisConn)
	SendToReplicas(*[]byte)
	UpdateReplica(string, int)
	ResetBytes()
}
