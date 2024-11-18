package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/mappers"
	"net"
)

func HandleError(conn *net.Conn, error error) {
	(*conn).Write([]byte(mappers.ErrorStringWithMessage(error)))
}