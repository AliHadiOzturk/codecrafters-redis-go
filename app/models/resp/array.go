package resp

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/models"
)

type RESPArray struct {
	models.Message
	models.MessageHandler
}

func NewArray(data []byte) *RESPArray {
	return &RESPArray{Message: models.Message{Data: data, Responses: map[string]interface{}{
		"PING": "+PONG",
	}}}
}

func (r *RESPArray) Process() []byte {
	r.Decode()
	r.Handle()
	return r.Response()
}

func (r *RESPArray) Decode() {
	stringData := string(r.Data)

	strs := strings.Split(stringData, "\r\n")

	if len(strs) < 3 {
		r.Command = ""
		return
	}

	message := strs[2]

	r.Command = message
}

func (r *RESPArray) Handle() {
	// Nothing to handle in simple string
}

func (r *RESPArray) Encode() []byte {
	response := r.Responses[r.Command]

	formatted := fmt.Sprintf("%s\r\n", response)

	return []byte(formatted)
}

func (r *RESPArray) Response() []byte {
	return r.Encode()
}
