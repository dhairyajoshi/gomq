package connections

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

type Message struct {
	FuncName string   `json:"func"`
	Args     []string `json:"args"`
}

func Connect(port string) {
	conn, err := net.Dial("tcp", port)
	if err != nil {
		fmt.Println("Error connecting to server: ", err.Error())
		return
	}
	reader := bufio.NewScanner(os.Stdin)

	for {
		buff := make([]byte, 1024)
		n, err := conn.Read(buff)
		if err != nil {
			fmt.Println("error reading from conn: ", err.Error())
			continue
		}
		fmt.Println(string(buff[:n]))
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
			continue
		}
		conn.Write(jsonData)
		if command == "exit" {
			return
		}
	}
}
