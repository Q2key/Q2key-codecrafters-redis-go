package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/app/core"
)

type Handler interface {
	Handle(core.RConn, []string, *[]byte)
}
