package utils

import "github.com/codecrafters-io/redis-starter-go/app/messages"

func MessageHandler(data []byte) messages.MessageHandler {
	if len(data) == 0 {
		return messages.NewSimpleString(data)
	}

	switch string(data[0:1]) {
	case "+":
		return messages.NewSimpleString(data)
	case "*":
		return messages.NewArray(data)
	// case "$":
	// 	return messages.NewArray(data)
	default:
		return messages.NewSimpleString(data)
	}
}
