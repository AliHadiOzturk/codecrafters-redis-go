package commands

import "fmt"

var information = map[string]string{
	"role": "master",
	// "connected_slaves": "",
}

func Info(parameters []string) (string, error) {

	var response string = ""
	var size int = 0
	for k, v := range information {

		response += k + ":" + v + "\r\n"

		size += len(k) + len(v) + 1
	}

	response = "$" + fmt.Sprint(size) + "\r\n" + response

	return response, nil
}
