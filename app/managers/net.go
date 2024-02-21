package managers

import (
	"fmt"
	"net"
	"os"
)

type NetManager struct {
	listener *net.Listener

	message chan []byte
}

func NewNetManager() *NetManager {
	return &NetManager{
		message: make(chan []byte),
	}
}

func (n *NetManager) Init() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")

	n.listener = &l

	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
}

func (n *NetManager) HandleConnection(connection net.Conn) {
	defer connection.Close()

	for {
		data := make([]byte, 1024)
	
		connection.Read(data)
	
		fmt.Println("Received data: ", string(data))
	
		_, err := connection.Write([]byte("+PONG\r\n"))
	
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
			os.Exit(1)
		}
	}

}

func (n *NetManager) Start() {
	for {
		connection, err := (*n.listener).Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
		}

		go n.HandleConnection(connection)
		// connection.Close()
	}
}

func (n *NetManager) Stop() {
	(*n.listener).Close()
}
