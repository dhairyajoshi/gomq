package main

import (
	"fmt"

	"github.com/dhairyajoshi/gomq/commands"
	"github.com/dhairyajoshi/gomq/connections"
)

func main() {
	port := ":8000"
	fmt.Println("Starting GOMQ on port: ", port)
	commands.RegisterCommands()
	connections.Listen(port)
}
