package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/floppyzedolfin/sync/replica/internal"
	pb "github.com/floppyzedolfin/sync/replica/replica"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 8405, "The cmd port")
	localReplica = flag.String("local_replica", "/data/replica", "Path to the local replica")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterReplicaServer(grpcServer, internal.NewServer(*localReplica))
	fmt.Println("starting replica server, replicating to " + *localReplica)
	grpcServer.Serve(lis)
}
