package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/adapters"
	"github.com/codecrafters-io/redis-starter-go/app/contracts"
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

func (h *PsyncHandler) Handle(conn *net.Conn, c contracts.Command[string]) {
	if c == nil || !c.Validate() {
		log.Fatal()
	}

	resp := adapters.FromString("FULLRESYNC 8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb 0\r\n")
	(*conn).Write([]byte(resp))
}
