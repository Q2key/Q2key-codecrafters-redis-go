package main

import (
	"fmt"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/repr"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func handleClient(conn net.Conn) {
	buff := make([]byte, 2024)

	for {
		conn.Read(buff)
		s := string(buff)

		//query without value type mark
		v := s[1:]
		c := repr.ParseArray(v)

		cmd := c[0]
		switch cmd {
		case "ECHO":
			conn.Write([]byte(fmt.Sprintf("+%s\r\n", c[1])))
		case "PING":
			conn.Write([]byte("+PONG\r\n"))
		}
	}

}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	//Uncomment this block to pass the first stage

	ln, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection")
			continue
		}

		go handleClient(conn)
	}
}
