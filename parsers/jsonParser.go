package parsers

import (
	"encoding/json"
	"fmt"
)

type JsonParser struct{}

func (jp JsonParser) Encode(response ServerResponse) []byte {
	data, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error encoding to json: ", err.Error())
		return []byte{}
	}
	return data
}

func (jp JsonParser) Decode(request []byte) DecodedMessage {
	jsonData := DecodedMessage{}
	err := json.Unmarshal(request, &jsonData)
	if err != nil {
		fmt.Println("Error decoding to json: ", err.Error())
		return DecodedMessage{}
	}
	return jsonData
}

func (jp JsonParser) ClientDecode(input []byte) ServerResponse {
	jsonData := ServerResponse{}
	err := json.Unmarshal(input, &jsonData)
	if err != nil {
		fmt.Println("Error decoding to json: ", err.Error())
		return ServerResponse{}
	}
	return jsonData
}

func NewJsonParser() JsonParser {
	return JsonParser{}
}
