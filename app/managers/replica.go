package managers

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/messages"
	"github.com/codecrafters-io/redis-starter-go/app/models"
)

type ReplicaManager struct {
	Replica    *models.Replica
	NetManager *NetManager
	Connection *net.TCPConn
}

func NewReplicaManager(netManager *NetManager, replica *models.Replica) *ReplicaManager {
	return &ReplicaManager{
		Replica:    replica,
		NetManager: netManager,
	}
}

func (rm *ReplicaManager) Init(port string) {

	if rm.Replica.IsMaster() {
		rm.NetManager.Listen(port)
		rm.NetManager.Start()
		return
	}

	conn, err := rm.NetManager.Connect(rm.Replica.MasterHost, rm.Replica.MasterPort)

	if err != nil {
		panic(err)
	}

	rm.Connection = conn

	_, err = rm.Connection.Write(messages.NewArray([]byte{}).Prepare("PING"))

	if err != nil {
		panic(err)
	}

	fmt.Printf("Connected to master at %s:%s\n", rm.Replica.MasterHost, rm.Replica.MasterPort)

	rm.Read()
}

func (rm *ReplicaManager) Start() {

}

func (rm *ReplicaManager) Read() {
	fmt.Println("Started reading from master...")
	for {
		read := make([]byte, 1024)

		count, err := rm.Connection.Read(read)

		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Master is down")
			}
		}

		if count == 0 {
			continue
		}

		message := read[:count]

		fmt.Println("Received data: ", string(message))
	}

}
