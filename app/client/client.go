package client

import (
	"bufio"
	"fmt"
	"net"
)

type TcpClient struct {
	host string
	port string
	addr string
	conn net.Conn
}

func NewTcpClient(host string, port string) *TcpClient {
	return &TcpClient{
		host: host,
		port: port,
		addr: host + ":" + port,
	}
}

func (r *TcpClient) Disconnect() {
	err := r.conn.Close()
	if err != nil {
		return
	}
}

func (r *TcpClient) UseConnection(conn net.Conn) {
	r.conn = conn
}

func (r *TcpClient) Connect() error {
	conn, err := net.Dial("tcp", r.addr)
	if err != nil {
		return err
	}

	r.conn = conn

	return nil
}

func (r *TcpClient) Conn() *net.Conn {
	return &r.conn
}

func (r *TcpClient) SendBytes(bytes []byte) (*[]byte, error) {
	//defer r.conn.Close()

	_, err := r.conn.Write(bytes)
	if err != nil {
		return nil, err
	}

	read := bufio.NewReader(r.conn)
	resp, err := read.ReadBytes('\n')

	if err != nil {
		fmt.Println("Error reading resp:", err)

	}

	return &resp, nil
}
