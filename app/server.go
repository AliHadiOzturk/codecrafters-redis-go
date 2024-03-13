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

	models.InitReplica(*replicaOf)

	clientManager := managers.NewClientManager()

	go clientManager.Init()

	netManager := managers.NewNetManager(clientManager)

	managers.NewReplicaManager(netManager, models.ReplicaInfo).Init(*port)

	// func() {
	// 	for {
	// 	}
	// }()
}
