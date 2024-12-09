package handlers

import (
	"encoding/json"
	"fmt"
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

	fmt.Println(args)
	rep, err := strconv.Atoi(args[1])
	if err != nil {
		return
	}

	cnt, err := strconv.Atoi(args[2])
	if err != nil {
		return
	}

	h.instance.ScheduleReplicas(rep, cnt)

	s := (*h).instance.GetScheduler()
	if (s) == nil {
		return
	}

	strCount := strconv.Itoa((*s).ActiveRepicasCount)

	b, err := json.MarshalIndent((s), "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(b))

	conn.Write([]byte(core.FromStringToRedisInteger(strCount)))
}
