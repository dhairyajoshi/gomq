package server

import (
	"fmt"
	"net"
)

func handleConn(conn net.Conn) {
	len, err := conn.Write([]byte("Connected to GoMQ, ready to exchange messages!\n>"))
	if err != nil {
		fmt.Println("error writing to connection, ", err.Error())
	}
	fmt.Println("wrote ", len, " bytes to connection")
	for {
		input := make([]byte, 1024)
		_, err := conn.Read(input)
		if err != nil {
			fmt.Println("Error reading input: ", err.Error())
			continue
		}
		fmt.Println("Received input: ", string(input))
		conn.Write([]byte("Received input!\n>"))
	}
}
