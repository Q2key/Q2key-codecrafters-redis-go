package core

import (
	"strconv"
	"time"
)

func NewWaitHandler(instance Redis) *WaitHandler {
	return &WaitHandler{
		instance: instance,
	}
}

type WaitHandler struct {
	instance Redis
}

func (h *WaitHandler) Handle(conn Conn, args []string, _ *[]byte) {
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

	for _, r := range *h.instance.RepConnPool {
		if r.Offset() >= bytesNeeded {
			done[r.Id()] = true
			continue
		}

		go func() {
			r.Conn().Write([]byte("*3\r\n$8\r\nREPLCONF\r\n$6\r\nGETACK\r\n$1\r\n*\r\n"))
		}()

	}

	ch := *h.instance.AckChan

awaitingLoop:
	for len(done) < rep {
		select {
		case c := <-ch:
			{
				if c.Offset >= bytesNeeded {
					done[c.ConnId] = true
				}
			}
		case <-time.After(time.Duration(timeout) * time.Millisecond):
			{
				break awaitingLoop
			}
		}
	}

	v := strconv.Itoa(len(done))
	conn.Conn().Write([]byte(FromStringToRedisInteger(v)))
}
