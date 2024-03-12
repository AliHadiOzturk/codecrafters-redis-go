package managers

import (
	"fmt"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/utils"
)

type Client struct {
	manager    *NetManager
	connection net.Conn
	send       chan []byte
	message    chan []byte
}

func NewClient(manager *NetManager, connection net.Conn) *Client {
	return &Client{
		manager:    manager,
		connection: connection,
		send:       make(chan []byte),
		message:    make(chan []byte),
	}
}

type NetManager struct {
	listener   *net.Listener
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewNetManager(port string) *NetManager {
	nm := NetManager{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}

	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))

	nm.listener = &l

	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	fmt.Printf("Listening on port %s\n", port)

	return &nm
}

func (c *Client) Init() {
	fmt.Println("Initializing client...")
	// defer c.manager.Stop()

	for {
		select {
		// case message := <-c.send:
		// 	_, err := c.connection.Write(message)

		// 	if err != nil {
		// 		fmt.Println("Error writing to connection: ", err.Error())
		// 		// c.manager.unregister <- c
		// 		// c.connection.Close()
		// 		// break
		// 	}
		case message := <-c.message:
			if len(message) == 0 {
				continue
			}
			response := utils.MessageHandler(message).Process()
			fmt.Println("Sending response: ", string(response))
			c.connection.Write(response)
		}
	}
}

func (c *Client) Read() {
	fmt.Println("Started reading from connection...")
	for {
		read := make([]byte, 1024)

		count, err := c.connection.Read(read)

		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Connection closed by client")
				c.manager.unregister <- c
				break
			}
		}

		if count == 0 {
			continue
		}

		message := read[:count]

		fmt.Println("Received data: ", string(message))
		c.message <- message
	}
}

func (n *NetManager) Init() {
	for {
		select {
		case client := <-n.register:
			n.clients[client] = true
			fmt.Println("Registered client...")
		case client := <-n.unregister:
			if _, ok := n.clients[client]; ok {
				delete(n.clients, client)
				close(client.send)
				close(client.message)
				// client.connection.Close()
				fmt.Println("Unregistered client...")
			}
		case message := <-n.broadcast:
			for client := range n.clients {
				client.send <- message
			}
		}
	}
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

		client := NewClient(n, connection)
		n.register <- client

		go client.Init()
		go client.Read()
	}
}

func (n *NetManager) Broadcast(message []byte) {
	for client := range n.clients {
		client.send <- message
	}
}

func (n *NetManager) Stop() {
	(*n.listener).Close()
}
