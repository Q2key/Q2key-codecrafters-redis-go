package contracts

import "net"

type Handler interface {
	Handle(net.Conn, Command)
}
