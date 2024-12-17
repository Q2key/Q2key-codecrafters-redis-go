package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/core/rconn"
)

type Handler interface {
	Handle(rconn.RConn, []string, *[]byte)
}
