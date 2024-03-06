package utils

import (
	"github.com/codecrafters-io/redis-starter-go/app/models"
	"github.com/codecrafters-io/redis-starter-go/app/models/resp"
)

func MessageHandler(data []byte) models.MessageHandler {
	if len(data) == 0 {
		return resp.NewSimpleString(data)
	}

	switch string(data[0:1]) {
	case "+":
		return resp.NewSimpleString(data)
	case "*":
		return resp.NewArray(data)
	// case "$":
	// 	return resp.NewArray(data)
	default:
		return resp.NewSimpleString(data)
	}
}
