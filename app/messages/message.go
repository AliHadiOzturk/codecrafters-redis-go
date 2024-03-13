package messages

type MessageHandler interface {
	Process() []byte
	Handle()
	Decode()
	Encode() []byte
	Response() []byte
	Prepare(command string) []byte
}

type Message struct {
	Data      []byte
	Request   string
	Args      []string
	Responses map[string]interface{}
	Commands  map[string]interface{}
	MessageHandler
}
