package parsers

type DecodedMessage struct {
	FuncName string `json:"func"`
	Args     []any  `json:"args"`
}

type Parser interface {
	Encode(map[string]any) []byte
	Decode(input []byte) DecodedMessage
}
