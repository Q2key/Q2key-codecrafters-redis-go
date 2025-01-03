package core

import (
	"encoding/hex"
	"fmt"
)

func NewPsyncHandler(instance Redis) *PsyncHandler {
	return &PsyncHandler{
		instance: instance,
	}
}

type PsyncHandler struct {
	instance Redis
}

func (h *PsyncHandler) Handle(conn Conn, _ []string, _ *[]byte) {

	h.instance.RegisterReplicaConn(conn)

	mess := fmt.Sprintf("FULLRESYNC %s 0", conn.Id())
	resp := FromStringToRedisCommonString(mess)

	rdb := "524544495330303131fa0972656469732d76657205372e322e30fa0a72656469732d62697473c040fa056374696d65c26d08bc65fa08757365642d6d656dc2b0c41000fa08616f662d62617365c000fff06e3bfec0ff5aa2"
	rdbBuff, _ := hex.DecodeString(rdb)

	chunkA := []byte(resp)
	chunkB := []byte("$88\r\n")
	chunkC := rdbBuff

	res := CombineBuffers(chunkA, chunkB, chunkC)

	conn.Conn().Write(res)
}

func CombineBuffers(buffs ...[]byte) []byte {
	buff := make([]byte, 0)
	for _, b := range buffs {
		buff = append(buff, b...)
	}
	return buff
}
