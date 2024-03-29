package messages

import (
	"fmt"
	"reflect"
	"strings"
)

type RESPBulkString struct {
	Message
	MessageHandler
}

func NewBulkString(data []byte) *RESPBulkString {
	return &RESPBulkString{Message: Message{Data: data,
		Responses: map[string]interface{}{},
		Commands: map[string]interface{}{
			"PING": "ping",
		}}}
}

func (r *RESPBulkString) Process() []byte {
	r.Decode()
	r.Handle()
	return r.Response()
}

func (r *RESPBulkString) Decode() {
	stringData := string(r.Data)

	strs := strings.Split(stringData, "\r\n")

	message := strs[1]

	r.Request = message
}

func (r *RESPBulkString) Handle() {
	// Nothing to handle in simple string
}

func (r *RESPBulkString) Encode() []byte {
	response := r.Responses[r.Request]

	formatted := fmt.Sprintf("+%s\r\n", response)

	return []byte(formatted)
}

func (r *RESPBulkString) Response() []byte {
	return r.Encode()
}

func (r *RESPBulkString) Prepare(command string) []byte {
	response := r.Commands[command]

	if response == nil {
		response = r.Commands[""]
	}

	if response != nil && reflect.TypeOf(response).Kind() == reflect.String {
		keywords := strings.Split(response.(string), " ")

		var newResponse string = ""

		for _, keyword := range keywords {
			newResponse = newResponse + fmt.Sprintf("$%d\r\n%s\r\n", len(keyword), keyword)
		}

		response = fmt.Sprintf("*%d\r\n%s", len(keywords), newResponse)
	}

	formatted := fmt.Sprintf("%s\r\n", response)

	return []byte(formatted)
}
