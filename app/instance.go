package main

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/repr"
)

type RedisInstance struct {
	Store map[string]string
}

func (ri *RedisInstance) Get(key string) string {
	return ri.Store[key]
}

func (ri *RedisInstance) Set(key string, value string) {
	ri.Store[key] = value
}

func (ri *RedisInstance) ToResponse(input string) string {
	return fmt.Sprintf("+%s\r\n", input)
}

func (ri *RedisInstance) handleClient(conn net.Conn) {
	buff := make([]byte, 2024)

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
			conn.Write([]byte(ri.ToResponse(val)))
		case "SET":
			key, val := c[1], c[2]
			ri.Set(key, val)
			conn.Write([]byte(ri.ToResponse("OK")))
		case "ECHO":
			conn.Write([]byte(ri.ToResponse(c[1])))
		case "PING":
			conn.Write([]byte("+PONG\r\n"))
		}
	}
}
