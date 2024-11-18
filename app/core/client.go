package core

import (
	"fmt"
	"net"
)

func ConnectToRedisInstance(host, port string) *net.Conn {
	address := fmt.Sprintf("%s:%s", host, port)
	con, err := net.Dial("tcp", address)
	if err != nil {
		panic(err)
	}

	return &con
}

func SendRedisMessage(conn *net.Conn, message string) {
	(*conn).Write([]byte(message))
}
