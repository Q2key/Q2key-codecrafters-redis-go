package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/core"
)

func NewInfoHandler(instance core.Redis) *InfoHandler {
	return &InfoHandler{
		instance: instance,
	}
}

type InfoHandler struct {
	instance core.Redis
}

func (h *InfoHandler) Handle(conn core.Conn, args []string, _ *[]byte) {
	r := h.instance.Config.GetReplica()
	res := "role:master"
	if r != nil {
		res = "role:slave"
	}

	for _, v := range args {
		if v == "replication" {
			res += ":master_replid:8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb"
			res += ":master_repl_offset:0"
		}
	}

	conn.Conn().Write([]byte(core.FromStringToRedisBulkString(res)))
}
