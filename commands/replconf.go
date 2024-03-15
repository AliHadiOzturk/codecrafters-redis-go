package commands

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/app/models"
)

const (
	// TODO: consider renaming
	ReplConfArgTypeCapa = "capa"
	RelConfArgTypePort  = "listening-port"
)

func ReplConfReceive(args []string) (string, error) {
	if len(args) < 2 {
		return "", models.NewNotEnoughArgsError("REPLCONF")
	}

	switch args[0] {
	case ReplConfArgTypeCapa:
		// implement psync2 and other capabilities

	case RelConfArgTypePort:
		models.ReplicaInfo.AddReplica(&models.Replica{
			Port: args[1],
		})

	default:
		return "", models.NewNotEnoughArgsError(fmt.Sprintf("ERR Unknown REPLCONF subcommand or wrong number of arguments for '%s' command", args[0]))
	}

	return "+OK", nil
}

func ReplConfSend(args []string) (string, error) {

	args = append([]string{"REPLCONF"}, args...)

	var newResponse string = ""

	for _, arg := range args {
		newResponse = newResponse + fmt.Sprintf("$%d\r\n%s\r\n", len(arg), arg)
	}

	response := fmt.Sprintf("*%d\r\n%s", len(args), newResponse)

	return response, nil
}
