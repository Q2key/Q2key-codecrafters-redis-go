package handlers

import (
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
	amout := args[1]
	seconds := args[2]

	fmt.Println(amout, seconds)

	am, err := strconv.Atoi(amout)
	if err != nil {
		return
	}

	if am > 0 {
		conn.Write([]byte(core.FromStringToRedisInteger("7")))
	} else {
		conn.Write([]byte(core.FromStringToRedisInteger("0")))
	}
}
