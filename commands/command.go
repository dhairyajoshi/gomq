package commands

import (
	"fmt"
	"net"

	"github.com/dhairyajoshi/gomq/parsers"
)

func getCommandAndArgs(input string) (func(args ...any) parsers.ServerResponse, []any) {
	parser := parsers.GetParser()
	decodedInput := parser.Decode([]byte(input))
	exec, exists := commandFactory[decodedInput.FuncName]
	if !exists {
		fmt.Println("no function named ", decodedInput.FuncName)
		return func(args ...any) parsers.ServerResponse {
			return parsers.ServerResponse{}
		}, []any{}
	}
	return exec, decodedInput.Args
}

func getResponse(conn *net.Conn, input string) parsers.ServerResponse {
	fun, args := getCommandAndArgs(input)
	funcArgs := []any{conn}
	funcArgs = append(funcArgs, args...)
	return fun(funcArgs...)
}

func HandleCommand(conn *net.Conn, input string) parsers.ServerResponse {
	return getResponse(conn, input)
}
