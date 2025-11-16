package connections

import (
	"fmt"
	"net"

	"github.com/dhairyajoshi/gomq/commands"
)

func handleConn(conn net.Conn) {
	_, err := conn.Write([]byte("Connected to GoMQ, ready to exchange messages!\n>"))
	if err != nil {
		fmt.Println("error writing to connection, ", err.Error())
		conn.Close()
		return
	}
	for {
		input := make([]byte, 1024)
		n, err := conn.Read(input)
		if err != nil {
			fmt.Println("Error reading input: ", err.Error())
			return
		}
		exit, response := commands.HandleCommand(string(input[:n]))
		server_response := fmt.Sprint(response, "\n>")
		conn.Write([]byte(server_response))
		if exit {
			conn.Close()
			return
		}
	}
}
