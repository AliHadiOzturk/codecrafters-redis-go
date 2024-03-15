package models

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/common"
)

var ReplicaInfo *Replica

type Replica struct {
	Port                string              `resp:"-"`
	Role                string              `resp:"role"`
	masterReplicaId     string              `resp:"master_replid"`
	masterReplicaOffset string              `resp:"master_repl_offset"`
	MasterHost          string              `resp:"-"`
	MasterPort          string              `resp:"-"`
	Replicas            map[string]*Replica `resp:"-"`
}

func InitReplica(port string, repllicaOf string) {
	ReplicaInfo = &Replica{
		Port:                port,
		Role:                "master",
		masterReplicaId:     common.GenerateRandom(40),
		masterReplicaOffset: "0",
	}

	if repllicaOf != "" {
		ReplicaInfo.Role = "slave"
		fmt.Println("Replicating to another server...", repllicaOf)
		ReplicaInfo.MasterHost = strings.Split(repllicaOf, ":")[0]
		ReplicaInfo.MasterPort = strings.Split(repllicaOf, ":")[1]
	}
}

func (s *Replica) IsMaster() bool {
	return s.Role == "master"
}

func (s *Replica) AddReplica(replica *Replica) {
	if s.Replicas == nil {
		s.Replicas = make(map[string]*Replica)
	}

	s.Replicas[replica.Port] = replica
}

func (s *Replica) RemoveReplica(port string) {
	delete(s.Replicas, port)
}
