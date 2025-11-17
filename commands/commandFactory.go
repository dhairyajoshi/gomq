package commands

import (
	"fmt"

	"github.com/dhairyajoshi/gomq/parsers"
)

var commandFactory = map[string]func(args ...any) parsers.ServerResponse{
	"exit": func(args ...any) parsers.ServerResponse {
		return parsers.ServerResponse{Data: "Closing connection, bye bye!", Type: "server_response", SendNext: false, Close: true}
	},
	"echo": func(args ...any) parsers.ServerResponse {
		return parsers.ServerResponse{Close: false, Data: fmt.Sprint(args...), SendNext: true, Type: "server_response"}
	},
}

func registerQueueCommands() {
	for k, v := range queueCommands {
		commandFactory[k] = v
	}
}

func RegisterCommands() {
	registerQueueCommands()
}
