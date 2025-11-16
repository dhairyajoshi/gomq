package main

import (
	"fmt"

	"github.com/dhairyajoshi/gomq/connections"
)

func main() {
	port := ":8000"
	fmt.Println("Starting GOMQ on port: ", port)
	connections.Listen(port)
}
