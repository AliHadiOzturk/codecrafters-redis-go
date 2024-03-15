package messages

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/models"
	"github.com/codecrafters-io/redis-starter-go/commands"
)

const (
	CommandTypePing        = "PING"
	CommandTypeEcho        = "ECHO"
	CommandTypeGet         = "GET"
	CommandTypeSet         = "SET"
	CommandTypeInfo        = "INFO"
	CommandTypeReplicaConf = "REPLCONF"
)

type RESPArray struct {
	Message
	MessageHandler
}

func NewArray(data []byte) *RESPArray {
	return &RESPArray{Message: Message{Data: data,
		Responses: map[string]interface{}{
			"":              "-ERR COMMAND NOT FOUND",
			CommandTypePing: "+PONG",
			CommandTypeEcho: commands.Echo,
			CommandTypeSet:  commands.Set,
			CommandTypeGet:  commands.Get,
			CommandTypeInfo: commands.Info(*models.ReplicaInfo),
			// Consider moving is to commands packages
			CommandTypeReplicaConf: commands.ReplConfReceive,
		},
		Commands: map[string]interface{}{
			CommandTypePing:        "ping",
			CommandTypeReplicaConf: commands.ReplConfSend,
		},
	}}
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
		r.Request = ""
		return
	}

	message := strs[2]

	r.Request = message
}

func (r *RESPArray) Handle() {

	strs := strings.Split(string(r.Data), "\r\n")

	parameterCount, _ := strconv.Atoi(strings.Split(string(strs[0]), "")[1])

	parameterCount = parameterCount - 1

	fmt.Printf("Parameter Count: %d\n", parameterCount)

	for i := 3; i < len(strs); i += 2 {
		if i+1 == len(strs) {
			break
		}
		r.Args = append(r.Args, strs[i+1])
	}

}

func (r *RESPArray) Encode() []byte {
	response := r.Responses[r.Request]

	if response == nil {
		response = r.Responses[strings.ToUpper(r.Request)]
	}

	if response != nil && reflect.TypeOf(response).Kind() == reflect.Func {
		resp, err := response.(func([]string) (string, error))(r.Args)

		if err != nil {
			response = fmt.Sprintf("-%s", err.Error())
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

func (r *RESPArray) Prepare(command string, args []string) []byte {
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

	if response != nil && reflect.TypeOf(response).Kind() == reflect.Func {
		resp, err := response.(func([]string) (string, error))(args)

		if err != nil {
			response = fmt.Sprintf("-%s", err.Error())
		} else {
			response = resp
		}
	}

	return []byte(response.(string))
}
