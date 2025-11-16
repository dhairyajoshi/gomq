package server

import (
	"fmt"
	"strings"
)

func handleCommand(input string) (bool, string) {
	fmt.Println(input)
	if strings.HasPrefix(input, "exit") {
		return true, "Closing connection, bye bye!"
	}
	switch input {
	case "exit":
		return true, "Closing connection, bye bye!"
	}
	return false, "response to request"
}
