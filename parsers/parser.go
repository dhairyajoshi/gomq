package parsers

type DecodedMessage struct {
	FuncName string   `json:"func"`
	Args     []string `json:"args"`
}

type Parser interface {
	Encode(map[string]any) []byte
	Decode(input []byte) DecodedMessage
}
