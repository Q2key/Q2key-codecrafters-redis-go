package redis

import (
	"net"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/repr"
)

type Handler struct {
	Store Store
}

func NewHandler() Handler {
	return Handler{
		Store: NewRedisStore(),
	}
}

type Command interface {
	Validate() bool
	Args() []string
	FromBytes(buff []byte) error
	Name() string
}

type GetCommand interface{}

func (h *Handler) HandleClient(conn net.Conn) {
	buff := make([]byte, 2024)

	ri := h.Store

	for {
		conn.Read(buff)
		s := string(buff)

		// query without value type mark
		v := s[1:]
		c := repr.ParseArray(v)

		cmd := c[0]
		switch cmd {
		case "CONFIG":
			if len(c) < 3 {
				conn.Write(repr.ToErrorString(nil))
				return
			}

			action, key := c[1], c[2]
			if action == "GET" && key == "dir" {
				resp := []string{key, h.Store.config.dir}
				conn.Write(repr.ToStringArray(resp))
				return
			}

			if action == "GET" && key == "dbfilename" {
				resp := []string{key, h.Store.config.dir}
				conn.Write(repr.ToStringArray(resp))
				return
			}

			conn.Write(repr.ToErrorString(nil))
		case "GET":
			key := c[1]
			val := ri.Get(key)
			if val.IsExpired() {
				conn.Write(repr.ToErrorString(nil))
			} else {
				conn.Write(repr.ToRegularString(val.Value))
			}
		case "SET":
			key, val := c[1], c[2]
			var exp int

			if len(c) >= 4 {
				exp, _ = strconv.Atoi(c[4])
			}
			ri.Set(key, val, int64(exp))

			conn.Write(repr.ToRegularString("OK"))
		case "ECHO":
			conn.Write(repr.ToRegularString(c[1]))
		case "PING":
			conn.Write(repr.ToRegularString("PONG"))
		}
	}
}
