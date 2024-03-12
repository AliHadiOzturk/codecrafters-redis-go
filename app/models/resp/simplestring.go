package resp

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/models"
)

type RESPSimpleString struct {
	models.Message
	models.MessageHandler
}

func NewSimpleString(data []byte) *RESPSimpleString {
	return &RESPSimpleString{Message: models.Message{Data: data, Request: "DEFAULT", Responses: map[string]interface{}{
		"":        "-ERR COMMAND NOT FOUND",
		"DEFAULT": "+OK",
		"PING":    "+PONG",
	}}}
}

func (r *RESPSimpleString) Process() []byte {
	r.Decode()
	r.Handle()
	return r.Response()
}

func (r *RESPSimpleString) Decode() {
	stringData := string(r.Data)

	strs := strings.Split(stringData, "\r\n")

	if len(strs) < 2 {
		r.Request = ""
		return
	} else if len(strs) == 2 {
		r.Request = strs[0]
		return
	}

	r.Request = strs[1]
}

func (r *RESPSimpleString) Handle() {
	// Nothing to handle in simple string

}

func (r *RESPSimpleString) Encode() []byte {
	response := r.Responses[r.Request]

	if response == nil {
		response = r.Responses[""]
	}

	formatted := fmt.Sprintf("%s\r\n", response)

	return []byte(formatted)
}

func (r *RESPSimpleString) Response() []byte {
	return r.Encode()
}
