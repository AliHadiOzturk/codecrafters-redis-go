package commands

import (
	"fmt"
	"reflect"

	"github.com/codecrafters-io/redis-starter-go/app/models"
)

func Info(parameters []string) (string, error) {

	var response string = ""
	var size int = 0
	var crlfCount int = 1

	values := reflect.ValueOf(*models.ServerInfo.ReplicaInfo)
	types := values.Type()
	for i := 0; i < values.NumField(); i++ {
		name, value := types.Field(i).Name, values.Field(i).String()
		response += name + ":" + value + "\r\n"
		crlfCount++
		size += len(name) + len(value) + 1
	}

	response = "$" + fmt.Sprint(size+crlfCount) + "\r\n" + response

	return response, nil
}
