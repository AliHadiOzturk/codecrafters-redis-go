package main

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/managers"
	"github.com/codecrafters-io/redis-starter-go/app/models"
)

func main() {

	port := flag.String("port", "6379", "Port to listen on")
	replicaHost := flag.String("replicaof", "", "Replicate to another server")
	flag.Parse()

	var replicaOf string

	args := flag.Args()
	if len(args) > 0 {
		lastArg := args[len(args)-1]
		port, err := strconv.Atoi(lastArg)
		if err != nil {
			fmt.Println("Last argument is not an integer")
		}

		if port > 0 && port < 65536 {
			replicaOf = fmt.Sprintf("%s:%d", *replicaHost, port)

		} else {
			fmt.Println("No non-flag arguments provided")
		}
	}

	models.InitReplica(replicaOf)

	clientManager := managers.NewClientManager()

	go clientManager.Init()

	netManager := managers.NewNetManager(clientManager)

	managers.NewReplicaManager(netManager, models.ReplicaInfo).Init(*port)

	// func() {
	// 	for {
	// 	}
	// }()
}
