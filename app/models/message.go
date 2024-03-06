package models

type MessageHandler interface {
	Process() []byte
	Handle()
	Decode()
	Encode() []byte
	Response() []byte
}

type Message struct {
	Data      []byte
	Command   string
	Args      []string
	Responses map[string]interface{}

	MessageHandler
}
