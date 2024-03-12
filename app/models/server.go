package models

import "github.com/codecrafters-io/redis-starter-go/app/common"

var ServerInfo *Server

func InitServer(repllicaOf string) {
	ServerInfo = &Server{
		ReplicaInfo: &Replication{
			role: func() string {
				if repllicaOf != "" {
					return "slave"
				} else {
					return "master"
				}
			}(),
			master_replid:      common.GenerateRandom(40),
			master_repl_offset: "0",
		},
	}
}

type Server struct {
	ReplicaInfo *Replication
}

type Replication struct {
	role               string
	master_replid      string
	master_repl_offset string
}
