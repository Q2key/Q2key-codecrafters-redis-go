package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/repr"
	"net"
)

func HandleError(conn *net.Conn, error error) {
	(*conn).Write([]byte(repr.ErrorStringWithMessage(error)))
}
