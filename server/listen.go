package server

import (
	"fmt"
	"net"
)

func Listen() {
	listner, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Println("error listening for connections!")
	}
	for {
		conn, err := listner.Accept()
		if err != nil {
			fmt.Println("error accepting connection!")
		}

		go handleConn(conn)
	}
}
