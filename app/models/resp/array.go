package resp

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/models"
	"github.com/codecrafters-io/redis-starter-go/commands"
)

type RESPArray struct {
	models.Message
	models.MessageHandler
}

func NewArray(data []byte) *RESPArray {
	return &RESPArray{Message: models.Message{Data: data, Responses: map[string]interface{}{
		"":     "-ERR COMMAND NOT FOUND",
		"PING": "+PONG",
		"ECHO": commands.Echo,
		"SET":  commands.Set,
		"GET":  commands.Get,
		"INFO": commands.Info,
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

	strs := strings.Split(string(r.Data), "\r\n")

	parameterCount, _ := strconv.Atoi(strings.Split(string(strs[0]), "")[1])

	parameterCount = parameterCount - 1

	println("Parameter Count: ", parameterCount)

	for i := 3; i < len(strs); i += 2 {
		if i+1 == len(strs) {
			break
		}
		r.Args = append(r.Args, strs[i+1])
	}

}

func (r *RESPArray) Encode() []byte {
	response := r.Responses[r.Command]

	if response == nil {
		response = r.Responses[strings.ToUpper(r.Command)]
	}

	if response != nil && reflect.TypeOf(response).Kind() == reflect.Func {
		resp, err := response.(func([]string) (string, error))(r.Args)

		if err != nil {
			// Implement error handling
		} else {
			response = resp
		}

	}

	if response == nil {
		response = r.Responses[""]
	}

	formatted := fmt.Sprintf("%s\r\n", response)

	return []byte(formatted)
}

func (r *RESPArray) Response() []byte {
	return r.Encode()
}
