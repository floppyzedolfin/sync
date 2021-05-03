package main

import (
	"context"
	"flag"
	"log"

	pb "github.com/floppyzedolfin/sync/replica/replica"
	"github.com/floppyzedolfin/sync/watcher/internal/watcher"
	"google.golang.org/grpc"
)

var (
	serverAddr = flag.String("server_addr", "localhost:1234", "The server address in the format of host:port")
	watchedDir = flag.String("watched_dir", "example/playground", "path from which changes will be propagated")
)

func main() {
	flag.Parse()

	// get a connection to the server
	conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("unable to dial a connection to server: %s", err.Error())
	}
	defer conn.Close()
	c := pb.NewReplicaClient(conn)
	if err != nil {
		log.Fatalf("unable to get a client to the replica server: %s", err.Error())
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	// start the watching process
	w, err := watcher.New(c, cancelFunc)
	if err != nil {
		log.Fatalf("error while instantiating the watcher: %s", err.Error())
	}
	w.Watch(ctx, *watchedDir)
}
