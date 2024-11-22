package handlers

import (
	"encoding/hex"
	"fmt"
	"github.com/codecrafters-io/redis-starter-go/app/contracts"
	"github.com/codecrafters-io/redis-starter-go/app/core"
	"log"
	"net"
)

func NewPsyncHandler(instance contracts.Instance) *PsyncHandler {
	return &PsyncHandler{
		instance: instance,
	}
}

type PsyncHandler struct {
	instance contracts.Instance
}

func (h *PsyncHandler) Handle(conn net.Conn, c contracts.Command) {
	if c == nil || !c.Validate() {
		log.Fatal()
	}

	mess := fmt.Sprintf("FULLRESYNC %s 0", h.instance.GetReplicaId())

	resp := core.FromStringToRedisCommonString(mess)

	rdb := "524544495330303131fa0972656469732d76657205372e322e30fa0a72656469732d62697473c040fa056374696d65c26d08bc65fa08757365642d6d656dc2b0c41000fa08616f662d62617365c000fff06e3bfec0ff5aa2"
	rdbBuff, _ := hex.DecodeString(rdb)

	conn.Write([]byte(resp))
	conn.Write([]byte("$88\r\n"))
	conn.Write([]byte(rdbBuff))
}
