package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/contracts"
	"github.com/codecrafters-io/redis-starter-go/app/core"
	"log"
)

func NewInfoHandler(instance contracts.Instance) *InfoHandler {
	return &InfoHandler{
		instance: instance,
	}
}

type InfoHandler struct {
	instance contracts.Instance
}

func (h *InfoHandler) Handle(conn contracts.RedisConn, c contracts.Command) {
	if c == nil || !c.Validate() {
		log.Fatal()
	}

	r := h.instance.GetConfig().GetReplica()
	res := "role:master"
	if r != nil {
		res = "role:slave"
	}

	args := c.Args()
	for _, v := range args {
		if v == "replication" {
			res += ":master_replid:8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb"
			res += ":master_repl_offset:0"
		}
	}

	conn.GetConn().Write([]byte(core.FromStringToRedisBulkString(res)))
}
