package handlers

import (
	"net"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/contracts"
	"github.com/codecrafters-io/redis-starter-go/app/core"
)

func NewWaitHandler(instance contracts.Instance) *WaitHandler {
	return &WaitHandler{
		instance: instance,
	}
}

type WaitHandler struct {
	instance contracts.Instance
}

func (h *WaitHandler) Handle(conn net.Conn, c contracts.Command) {
	if c == nil || !c.Validate() {
		return
	}

	args := c.Args()
	_ = args[1]
	_ = args[2]

	rlen := len(h.instance.GetReplicas())
	slen := strconv.Itoa(rlen)

	conn.Write([]byte(core.FromStringToRedisInteger(slen)))
}
