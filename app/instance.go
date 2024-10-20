package main

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/repr"
)

type RedisInstance struct {
	Store map[string]RedisValue
}

type RedisValue struct {
	Value   string
	Expired time.Time
	Created time.Time
}

func (rv *RedisValue) IsExpired() bool {
	if rv.Expired.IsZero() {
		return false
	}

	now := time.Now()

	diff := rv.Expired.UnixMilli() - now.UnixMilli()

	return diff <= 0
}

func (ri *RedisInstance) Get(key string) RedisValue {
	return ri.Store[key]
}

func (ri *RedisInstance) Set(key string, value string, expired int64) {
	now := time.Now().UTC()
	var exp time.Time
	if expired != 0 {
		exp = now.Add(time.Duration(expired) * time.Millisecond)
	}

	ri.Store[key] = RedisValue{Value: value, Created: now, Expired: exp}
}

func (ri *RedisInstance) ToOkString(input string) []byte {
	return []byte(fmt.Sprintf("+%s\r\n", input))
}

func (ri *RedisInstance) ToErrorString() []byte {
	return []byte("$-1\r\n")
}

func (ri *RedisInstance) HandleClient(conn net.Conn) {
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
