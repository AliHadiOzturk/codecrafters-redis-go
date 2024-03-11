package utils

// RESP: Redis serialization protocol
type RESP struct {
}

type SimpleString struct {
	resp RESP
}

func decode(data []byte) []byte {
	return data
}
