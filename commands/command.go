package commands

import (
	"fmt"

	"github.com/dhairyajoshi/gomq/io"
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

func getResponse(ioHandler *io.IOHandler, input string) parsers.ServerResponse {
	fun, args := getCommandAndArgs(input)
	funcArgs := []any{ioHandler}
	funcArgs = append(funcArgs, args...)
	return fun(funcArgs...)
}

func HandleCommand(ioHandler *io.IOHandler, input string) parsers.ServerResponse {
	return getResponse(ioHandler, input)
}
