package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/core"
	"net"
)

func HandleError(conn *net.Conn, error error) {
	(*conn).Write([]byte(core.ToRedisErrorStringWithMessage(error)))
}
