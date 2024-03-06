package resp

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/models"
)

type RESPBulkString struct {
	models.Message
	models.MessageHandler
}

func NewBulkString(data []byte) *RESPBulkString {
	return &RESPBulkString{Message: models.Message{Data: data, Responses: map[string]interface{}{}}}
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

	r.Command = message
}

func (r *RESPBulkString) Handle() {
	// Nothing to handle in simple string
}

func (r *RESPBulkString) Encode() []byte {
	response := r.Responses[r.Command]

	formatted := fmt.Sprintf("+%s\r\n", response)

	return []byte(formatted)
}

func (r *RESPBulkString) Response() []byte {
	return r.Encode()
}
