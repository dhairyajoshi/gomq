package connections

import (
	"fmt"
	"net"

	"github.com/dhairyajoshi/gomq/commands"
	"github.com/dhairyajoshi/gomq/io"
	"github.com/dhairyajoshi/gomq/parsers"
)

func handleConn(conn net.Conn) {
	ioHandler := io.NewIoHandler(&conn)
	_, err := ioHandler.Write(parsers.ServerResponse{Type: "server_response", Data: "Connected to GoMQ, ready to exchange messages!\n", SendNext: true, Close: false})
	if err != nil {
		fmt.Println("error writing to connection, ", err.Error())
		conn.Close()
		return
	}
	for {
		input, err := ioHandler.Read()
		if err != nil {
			fmt.Println("Error reading input: ", err.Error())
			return
		}
		response := commands.HandleCommand(ioHandler, string(input))
		ioHandler.Write(response)
		if response.Close {
			conn.Close()
			return
		}
	}
}
