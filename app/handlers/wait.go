package handlers

import (
	"fmt"
	"strconv"
	"time"

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

func (h *WaitHandler) Handle(conn contracts.RedisConn, c contracts.Command) {
	if c == nil || !c.Validate() {
		return
	}

	args := c.Args()

	rep, err := strconv.Atoi(args[1])
	if err != nil {
		return
	}

	timeout, err := strconv.Atoi(args[2])
	if err != nil {
		return
	}

	now := time.Now()
	til := now.Add(time.Duration(timeout) * time.Millisecond)

	for _, r := range h.instance.GetReplicas() {
		go func() {
			r.GetConn().Write([]byte("*3\r\n$8\r\nREPLCONF\r\n$6\r\nGETACK\r\n$1\r\n*\r\n"))
		}()
	}

	done := map[string]bool{}
	for len(done) < rep {
		ch := <-(*h.instance.GetAckChan())
		done[ch.ConnId] = true
		if now.UnixNano() > til.UnixNano() {
			fmt.Printf("\r\nreplica acked: %s %d", ch.ConnId, ch.Offset)
			break
		}
	}

	v := strconv.Itoa(len(done))
	conn.GetConn().Write([]byte(core.FromStringToRedisInteger(v)))
}
