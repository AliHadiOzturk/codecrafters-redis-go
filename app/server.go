package main

import (
	"github.com/codecrafters-io/redis-starter-go/app/managers"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	// fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	netManager := managers.NewNetManager()
	netManager.Init()
	netManager.Start()

}
