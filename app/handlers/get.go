package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/contracts"
	"github.com/codecrafters-io/redis-starter-go/app/mappers"
	"log"
	"net"
)

func NewGetHandler(instance contracts.Instance) *GetHandler {
	return &GetHandler{
		instance: instance,
	}
}

type GetHandler struct {
	instance contracts.Instance
}

func (h *GetHandler) Handle(conn *net.Conn, c contracts.Command[string]) {
	if c == nil || !c.Validate() {
		log.Fatal()
	}

	key := c.Args()[1]
	val := (h.instance).Get(key)

	if val.IsExpired() {
		(*conn).Write([]byte(mappers.ErrorString()))
	} else {
		(*conn).Write([]byte(mappers.FromString(val.GetValue())))
	}
}