package handlers

import (
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
	conn.Write([]byte(resp))
}
