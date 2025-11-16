package connections

import (
	"fmt"
	"net"
)

func Listen(port string) {
	listner, err := net.Listen("tcp", port)
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
