package connections

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/dhairyajoshi/gomq/io"
	"github.com/dhairyajoshi/gomq/messages"
	"github.com/dhairyajoshi/gomq/parsers"
)

var parser = parsers.GetParser()

type Message struct {
	FuncName string   `json:"func"`
	Args     []string `json:"args"`
}

func typeAssertedVal[T any](val any, def T) T {
	assertedVal, ok := val.(T)
	if !ok {
		fmt.Println("couldn't convert value to type: ", val)
		return def
	}
	return assertedVal
}

func handleResp(resp []byte) bool {
	serverResponse := parser.ClientDecode(resp)
	switch serverResponse.Type {
	case "server_response":
		fmt.Println(serverResponse.Data)
	case "message":
		responseData := typeAssertedVal(serverResponse.Data, map[string]any{})
		messageId := typeAssertedVal(responseData["id"], "")
		messageContent := typeAssertedVal(responseData["data"], "")
		message := messages.Message{Id: messageId, Data: messageContent}
		fmt.Println(message)
	}
	return serverResponse.SendNext
}

func sendCommand(conn *net.Conn, reader *bufio.Scanner) {
	fmt.Print("Enter command: ")
	reader.Scan()
	command := reader.Text()
	fmt.Print("Enter args (space separated): ")
	reader.Scan()
	rawArgs := reader.Text()
	args := strings.Fields(rawArgs) // split by spaces
	jsonData, err := json.Marshal(Message{FuncName: command, Args: args})
	if err != nil {
		fmt.Println("error converting command: ", err.Error())
		return
	}
	(*conn).Write(jsonData)
	if command == "exit" {
		return
	}
}

func Connect(port string) {
	conn, err := net.Dial("tcp", port)
	if err != nil {
		fmt.Println("Error connecting to server: ", err.Error())
		return
	}
	ioHandler := io.NewIoHandler(&conn)
	reader := bufio.NewScanner(os.Stdin)

	for {
		resp, err := ioHandler.Read()
		if err != nil {
			fmt.Println("error reading from conn: ", err.Error())
			continue
		}
		nextInput := handleResp(resp)
		if nextInput {
			sendCommand(&conn, reader)
		}
	}
}
