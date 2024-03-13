package models

import (
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/common"
)

var ReplicaInfo *Replica

type Replica struct {
	Role                string `resp:"role"`
	masterReplicaId     string `resp:"master_replid"`
	masterReplicaOffset string `resp:"master_repl_offset"`
	MasterHost          string `resp:"-"`
	MasterPort          string `resp:"-"`
}

func InitReplica(repllicaOf string) {
	ReplicaInfo = &Replica{
		Role:                "master",
		masterReplicaId:     common.GenerateRandom(40),
		masterReplicaOffset: "0",
	}

	if repllicaOf != "" {
		ReplicaInfo.Role = "slave"
		ReplicaInfo.MasterHost = strings.Split(repllicaOf, " ")[0]
		ReplicaInfo.MasterPort = strings.Split(repllicaOf, " ")[1]
	}
}

func (s *Replica) IsMaster() bool {
	return s.Role == "master"
}
