package managers

import (
	"fmt"
	"net"
	"os"
)

type NetManager struct {
	listener      *net.Listener
	clientManager *ClientManager
	connection    *net.TCPConn
}

func NewNetManager(clientManager *ClientManager) *NetManager {
	return &NetManager{
		clientManager: clientManager,
	}
}

func (n *NetManager) Listen(port string) {
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))

	n.listener = &l

	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	fmt.Printf("Listening on port %s\n", port)
}

func (n *NetManager) Connect(host string, port string) (*net.TCPConn, error) {
	addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%s", host, port))

	conn, err := net.DialTCP("tcp", nil, addr)

	if err != nil {
		fmt.Printf("Failed to connect")
		return nil, err
	}

	return conn, nil
}

func (n *NetManager) Start() {
	fmt.Println("Starting server...")
	for {
		connection, err := (*n.listener).Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			(*n.listener).Close()
			continue
		}

		fmt.Println("Accepted connection from: ", connection.RemoteAddr().String())

		client := NewClient(n.clientManager, connection)
		n.clientManager.register <- client

		go client.Init()
		go client.Read()
	}
}

func (n *NetManager) Broadcast(message []byte) {
	for client := range n.clientManager.clients {
		client.send <- message
	}
}

func (n *NetManager) Stop() {
	(*n.listener).Close()
}
