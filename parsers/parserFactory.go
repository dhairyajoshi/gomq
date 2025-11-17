package parsers

import "os"

func GetParser() Parser {
	protocol, found := os.LookupEnv("protocol")
	if !found {
		protocol = "json"
	}
	switch protocol {
	case "json":
		return NewJsonParser()
	}
	return NewJsonParser()
}
