package client

import (
	"net"
)

type TcpClient struct {
	host string
	port string
	addr string
	Conn *net.Conn
}

func NewTcpClient(host string, port string) *TcpClient {
	return &TcpClient{
		host: host,
		port: port,
		addr: host + ":" + port,
	}
}

func (r *TcpClient) Disconnect() {
	err := (*r.Conn).Close()
	if err != nil {
		return
	}
}

func (r *TcpClient) Connect() {
	conn, err := net.Dial("tcp", r.addr)
	if err != nil {
		return
	}

	r.Conn = &conn
}
