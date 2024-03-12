package models

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
		},
	}
}

type Server struct {
	ReplicaInfo *Replication
}

type Replication struct {
	role string
}
