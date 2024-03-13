package managers

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/utils"
)

type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (manager *ClientManager) Init() {
	fmt.Println("Initializing client manager...")
	for {
		select {
		case client := <-manager.register:
			manager.clients[client] = true
			fmt.Println("Registered client...")
		case client := <-manager.unregister:
			if _, ok := manager.clients[client]; ok {
				delete(manager.clients, client)
				close(client.send)
				close(client.message)
				// client.connection.Close()
				fmt.Println("Unregistered client...")
			}
		case message := <-manager.broadcast:
			for client := range manager.clients {
				client.send <- message
			}
		}
	}
}

type Client struct {
	manager    *ClientManager
	connection net.Conn
	send       chan []byte
	message    chan []byte
}

func NewClient(manager *ClientManager, connection net.Conn) *Client {
	return &Client{
		manager:    manager,
		connection: connection,
		send:       make(chan []byte),
		message:    make(chan []byte),
	}
}

func (c *Client) Init() {
	fmt.Println("Initializing client...")
	// defer c.manager.Stop()

	for {
		select {

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
