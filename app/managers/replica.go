package managers

import (
	"fmt"
	"net"
	"time"

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

	// if rm.Replica.IsMaster() {
	rm.NetManager.Listen(port)
	// return
	// }

	if !rm.Replica.IsMaster() {
		go rm.NetManager.Start()
		for {

			conn, err := rm.NetManager.Connect(rm.Replica.MasterHost, rm.Replica.MasterPort)

			if err != nil {
				fmt.Printf("Failed to connect to master at %s:%s", rm.Replica.MasterHost, rm.Replica.MasterPort)
				time.Sleep(10 * time.Second)
				continue
			}

			rm.Connection = conn

			_, err = rm.Connection.Write(messages.NewArray([]byte{}).Prepare(messages.CommandTypePing, []string{}))

			if err != nil {
				fmt.Println("Failed to write to master")
				time.Sleep(10 * time.Second)
				continue
			}
			
			// TODO: Workaround for the time being
			time.Sleep(100 * time.Millisecond)

			_, err = rm.Connection.Write(messages.NewArray([]byte{}).Prepare(messages.CommandTypeReplicaConf, []string{"listening-port", port}))

			if err != nil {
				fmt.Println("Failed to write to master")
				time.Sleep(10 * time.Second)
				continue
			}

			// TODO: Workaround for the time being
			time.Sleep(100 * time.Millisecond)

			_, err = rm.Connection.Write(messages.NewArray([]byte{}).Prepare(messages.CommandTypeReplicaConf, []string{"capa", "psync2"}))

			if err != nil {
				fmt.Println("Failed to write to master")
				time.Sleep(10 * time.Second)
				continue
			}

			fmt.Printf("Connected to master at %s:%s\n", rm.Replica.MasterHost, rm.Replica.MasterPort)

			rm.Read()
		}
	} else {
		rm.NetManager.Start()
	}
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
				break
			}
		}

		if count == 0 {
			continue
		}

		message := read[:count]

		fmt.Println("Received data: ", string(message))
	}

}
