package commands

import (
	"fmt"
	"reflect"

	"github.com/codecrafters-io/redis-starter-go/app/models"
)

func Info(parameters []string) (string, error) {

	var response string = ""
	var size int = 0

	values := reflect.ValueOf(*models.ServerInfo.ReplicaInfo)
	types := values.Type()
	for i := 0; i < values.NumField(); i++ {
		name, value := types.Field(i).Name, values.Field(i)
		response += name + ":" + value.String() + "\r\n"
		size += len(name) + len(value.String()) + 1
	}

	response = "$" + fmt.Sprint(size) + "\r\n" + response

	return response, nil
}
