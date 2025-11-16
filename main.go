package main

import (
	"fmt"
	"os"

	"github.com/dhairyajoshi/gomq/commands"
	"github.com/dhairyajoshi/gomq/connections"
	"github.com/dhairyajoshi/gomq/queues"
)

func main() {
	port := ":8000"
	commands.RegisterCommands()
	mode := os.Args
	if mode[1] == "serve" {
		fmt.Println("Starting GOMQ on port: ", port)
		go queues.MonitorQueues()
		connections.Listen(port)
	} else {
		connections.Connect(port)
	}
}
