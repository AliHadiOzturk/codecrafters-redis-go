package managers

import (
	"fmt"
	"net"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/messages"
	"github.com/codecrafters-io/redis-starter-go/app/models"
	"github.com/codecrafters-io/redis-starter-go/commands"
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
		// for {

		conn, err := rm.NetManager.Connect(rm.Replica.MasterHost, rm.Replica.MasterPort)

		if err != nil {
			fmt.Printf("Failed to connect to master at %s:%s", rm.Replica.MasterHost, rm.Replica.MasterPort)
			time.Sleep(10 * time.Second)
			// continue
		}

		rm.Connection = conn

		err = rm.SendPingToMaster()

		if err != nil {
			fmt.Println("Failed to send ping to master")
			time.Sleep(10 * time.Second)
			// continue
		}

		err = rm.SendReplConf(commands.RelConfArgTypePort)

		if err != nil {
			fmt.Println("Failed to send listening port to master")
			time.Sleep(10 * time.Second)
			// continue
		}

		err = rm.SendReplConf(commands.ReplConfArgTypeCapa)

		if err != nil {
			fmt.Println("Failed to send capa to master")
			time.Sleep(10 * time.Second)
			// continue
		}

		fmt.Printf("Connected to master at %s:%s\n", rm.Replica.MasterHost, rm.Replica.MasterPort)

		rm.Read()
		// }
	} else {
		rm.NetManager.Start()
	}
}

func (rm *ReplicaManager) SendPingToMaster() error {
	buff := make([]byte, 100)
	_, err := rm.Connection.Write(messages.NewArray([]byte{}).Prepare(messages.CommandTypePing, []string{}))

	if err != nil {
		return err
	}

	_, err = rm.Connection.Read(buff)

	if err != nil {
		return err
	}

	fmt.Println("Received response from master for ping: ", string(buff))

	return nil
}

func (rm *ReplicaManager) SendReplConf(commandType string) error {
	buff := make([]byte, 100)

	var err error

	switch commandType {
	case commands.RelConfArgTypePort:
		_, err = rm.Connection.Write(messages.NewArray([]byte{}).Prepare(messages.CommandTypeReplicaConf, []string{"listening-port", rm.Replica.Port}))
	case commands.ReplConfArgTypeCapa:
		_, err = rm.Connection.Write(messages.NewArray([]byte{}).Prepare(messages.CommandTypeReplicaConf, []string{"capa", "psync2"}))
	}

	if err != nil {
		return err
	}

	_, err = rm.Connection.Read(buff)

	if err != nil {
		return err
	}

	fmt.Printf("Received response from master for ReplConf(%s): %s", commandType, string(buff))

	return nil
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
