package client

import (
	"context"
	"fmt"
	"time"

	reference2 "github.com/floppyzedolfin/sync/reference/client/reference"
	"google.golang.org/grpc"
)

//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative reference/reference.proto

const (
	// default values, hardcoded for this exercise
	ip         = "172.18.0.23"
	portNumber = 8405
)

// Notify exposes the service's Notify endpoint
func Notify(ctx context.Context, req *reference2.NotifyRequest) (*reference2.NotifyResponse, error) {
	clientConnection, err := grpc.Dial(fmt.Sprintf("%s:%d", ip, portNumber), grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("unable to craete connection to reference: %w", err)
	}
	defer clientConnection.Close()

	client := reference2.NewReferenceClient(clientConnection)
	// setting a 1 sec timeout for all the requests.
	// Using a "local" network makes file transfer easy - this time could, say, be based on the volume of the request.
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	// handle timeouts properly
	defer cancel()

	// call the endpoint
	res, err := client.Notify(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error while querying the notify endpoint: %w", err)
	}

	return res, nil
}
