package redis

import (
	"github.com/codecrafters-io/redis-starter-go/app/repr"
	"net"
	"strconv"
)

type Handler struct {
	Store Store
}

func NewHandler() Handler {
	return Handler{
		Store: NewRedisStore(),
	}
}

func (h *Handler) HandleClient(conn net.Conn) {
	buff := make([]byte, 2024)

	ri := h.Store

	for {
		conn.Read(buff)
		s := string(buff)

		//query without value type mark
		v := s[1:]
		c := repr.ParseArray(v)

		cmd := c[0]
		switch cmd {
		case "GET":
			key := c[1]
			val := ri.Get(key)
			if val.IsExpired() {
				conn.Write(ri.ToErrorString())
			} else {
				conn.Write(ri.ToOkString(val.Value))
			}
		case "SET":
			key, val := c[1], c[2]
			var exp int

			if len(c) >= 4 {
				exp, _ = strconv.Atoi(c[4])
			}

			ri.Set(key, val, int64(exp))

			conn.Write(ri.ToOkString("OK"))
		case "ECHO":
			conn.Write(ri.ToOkString(c[1]))
		case "PING":
			conn.Write(ri.ToOkString("PONG"))
		}
	}
}
