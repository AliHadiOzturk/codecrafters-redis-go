package utils

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/messages"
)

type ReplicaUtils struct {
}


func (r *ReplicaUtils) Ping(connection net.Conn) {
	_, err := connection.Write(messages.NewArray([]byte{}).Prepare("PING"))

	if err != nil {
		panic("Failed to ping master")
	}
}
