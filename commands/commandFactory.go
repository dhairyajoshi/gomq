package commands

import "fmt"

var commandFactory = map[string]func(args ...any) (bool, string){
	"exit": func(args ...any) (bool, string) {
		return true, "Closing connection, bye bye!"
	},
	"echo": func(args ...any) (bool, string) {
		return false, fmt.Sprint(args...)
	},
}
