package parsers

type DecodedMessage struct {
	FuncName string `json:"func"`
	Args     []any  `json:"args"`
}

type ServerResponse struct {
	Data     any    `json:"data"`      // actual response data
	Type     string `json:"type"`      // 'server_response', 'message'
	SendNext bool   `json:"send_next"` // true to continue sending next commands, false to keep listening for response
	Close    bool   `json:"close"`
}

type Parser interface {
	Encode(response ServerResponse) []byte
	Decode(input []byte) DecodedMessage
	ClientDecode([]byte) ServerResponse
}
