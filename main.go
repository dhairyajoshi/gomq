package main

import (
	"fmt"

	"github.com/dhairyajoshi/gomq/server"
)

func main() {
	port := 8000
	fmt.Println("Starting GOMQ on port: ", port)
	server.Listen()
}
