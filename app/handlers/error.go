package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/adapters"
	"net"
)

func HandleError(conn *net.Conn, error error) {
	(*conn).Write([]byte(adapters.ErrorStringWithMessage(error)))
}
