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

	rep, err := strconv.Atoi(args[1])
	if err != nil {
		return
	}

	cnt, err := strconv.Atoi(args[2])
	if err != nil {
		return
	}

	h.instance.ScheduleReplicas(rep, cnt)

	var repc int

	shd := (*h).instance.GetScheduler()
	if shd != nil {
		repc = (*shd).ActiveRepicasCount
	} else {
		repc = len(h.instance.GetReplicas())
	}

	core.TraceObj(shd, "wait handler: ", core.Gray)

	strCount := strconv.Itoa(repc)
	conn.Write([]byte(core.FromStringToRedisInteger(strCount)))
}
