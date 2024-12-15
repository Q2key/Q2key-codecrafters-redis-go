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

	done := map[string]bool{}
	bytesNeeded := h.instance.GetWrittenBytes()

	for _, r := range h.instance.GetReplicas() {
		// reader := bufio.NewReader(r.GetConn())
		fmt.Println(r.GetOffset(), bytesNeeded)
		if r.GetOffset() >= bytesNeeded {
			done[r.GetId()] = true
			continue
		}

		go func() {
			fmt.Println("sending...", r.GetId(), r.GetOffset())
			r.GetConn().Write([]byte("*3\r\n$8\r\nREPLCONF\r\n$6\r\nGETACK\r\n$1\r\n*\r\n"))
		}()

	}

	ch := (*h.instance.GetAckChan())
awaitingLoop:
	for len(done) < rep {
		select {
		case c := <-ch:
			{
				fmt.Println("ACKED", c.Offset, bytesNeeded)

				if c.Offset >= bytesNeeded {
					done[c.ConnId] = true
				}
			}
		case t := <-time.After(time.Duration(timeout) * time.Millisecond):
			{
				fmt.Println("BREAKING LOOP", t.UTC())
				break awaitingLoop
			}
		}
	}

	v := strconv.Itoa(len(done))
	conn.GetConn().Write([]byte(core.FromStringToRedisInteger(v)))
}
