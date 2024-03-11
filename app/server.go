package main

import (
	"flag"

	"github.com/codecrafters-io/redis-starter-go/app/managers"
)

func main() {

	port := flag.String("port", "6379", "Port to listen on")
	flag.Parse()

	netManager := managers.NewNetManager(*port)
	go netManager.Init()
	netManager.Start()

	// func() {
	// 	for {
	// 	}
	// }()
}
