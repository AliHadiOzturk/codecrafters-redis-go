package main

import (
	"flag"

	"github.com/codecrafters-io/redis-starter-go/app/managers"
	"github.com/codecrafters-io/redis-starter-go/app/models"
)

func main() {

	port := flag.String("port", "6379", "Port to listen on")
	replicaOf := flag.String("replicaof", "", "Replicate to another server")
	flag.Parse()

	models.InitServer(*replicaOf)

	netManager := managers.NewNetManager(*port)
	go netManager.Init()
	netManager.Start()

	// func() {
	// 	for {
	// 	}
	// }()
}
