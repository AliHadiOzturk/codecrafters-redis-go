package commands

import (
	"fmt"
	"reflect"
)

func Info(replica interface{}) func(parameters []string) (string, error) {
	return func(parameters []string) (string, error) {
		var response string = ""
		var size int = 0
		var crlfCount int = 1

		values := reflect.ValueOf(replica)
		types := values.Type()
		for i := 0; i < values.NumField(); i++ {
			name, value := types.Field(i).Name, values.Field(i).String()

			tag := types.Field(i).Tag.Get("resp")

			if tag == "-" {
				continue
			}

			name = tag

			response += name + ":" + value + "\r\n"
			crlfCount++
			size += len(name) + len(value) + 1
		}

		response = "$" + fmt.Sprint(size+crlfCount) + "\r\n" + response

		return response, nil
	}
}
