package connections

import (
	"fmt"
	"net"

	"github.com/dhairyajoshi/gomq/commands"
	"github.com/dhairyajoshi/gomq/parsers"
)

func handleConn(conn net.Conn) {
	parser := parsers.GetParser()
	response := parser.Encode(parsers.ServerResponse{Type: "server_response", Data: "Connected to GoMQ, ready to exchange messages!\n", SendNext: true, Close: false})
	_, err := conn.Write(response)
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
		response := commands.HandleCommand(&conn, string(input[:n]))
		server_response := parser.Encode(response)
		conn.Write(server_response)
		if response.Close {
			conn.Close()
			return
		}
	}
}
