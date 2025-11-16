package commands

import (
	"fmt"
	"os"

	"github.com/dhairyajoshi/gomq/parsers"
)

func getParser() parsers.Parser {
	protocol, found := os.LookupEnv("protocol")
	if !found {
		protocol = "json"
	}
	switch protocol {
	case "json":
		return parsers.NewJsonParser()
	}
	return parsers.NewJsonParser()
}

func getCommandAndArgs(input string) (func(args ...any) (bool, string), []any) {
	parser := getParser()
	decodedInput := parser.Decode([]byte(input))
	exec, exists := commandFactory[decodedInput.FuncName]
	if !exists {
		fmt.Println("no function named ", decodedInput.FuncName)
		return func(args ...any) (bool, string) {
			return false, ""
		}, []any{}
	}
	return exec, decodedInput.Args
}

func getResponse(input string) (bool, string) {
	fun, args := getCommandAndArgs(input)
	return fun(args...)
}

func HandleCommand(input string) (bool, string) {
	return getResponse(input)
}
