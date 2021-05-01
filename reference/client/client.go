package client

import (
	"context"
	"fmt"
	"time"

	"github.com/floppyzedolfin/sync/reference/client/reference"
	"google.golang.org/grpc"
)

//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative reference/reference.proto

const (
	// default values, hardcoded for this exercise
	ip         = "172.18.0.23"
	PortNumber = 8405
)

// Patch exposes the service's Patch endpoint
func (c) Patch(ctx context.Context, req *reference.PatchRequest) (*reference.PatchResponse, error) {
	fmt.Printf("received patch request")
	clientConnection, err := grpc.Dial(fmt.Sprintf("%s:%d", ip, PortNumber), grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("unable to craete connection to reference: %w", err)
	}
	defer clientConnection.Close()

	client := reference.NewReferenceClient(clientConnection)
	// setting a 1 sec timeout for all the requests.
	// Using a "local" network makes file transfer easy - this time could, say, be based on the volume of the request.
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	// handle timeouts properly
	defer cancel()

	// call the endpoint
	res, err := client.Patch(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error while querying the patch endpoint: %w", err)
	}

	return res, nil
}

// Delete exposes the service's Delete endpoint
func (c) Delete(ctx context.Context, req *reference.DeleteRequest) (*reference.DeleteResponse, error) {
	clientConnection, err := grpc.Dial(fmt.Sprintf("%s:%d", ip, PortNumber), grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("unable to craete connection to reference: %w", err)
	}
	defer clientConnection.Close()

	client := reference.NewReferenceClient(clientConnection)
	// setting a 1 sec timeout for all the requests.
	// Using a "local" network makes file transfer easy - this time could, say, be based on the volume of the request.
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	// handle timeouts properly
	defer cancel()

	// call the endpoint
	res, err := client.Delete(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error while querying the delete endpoint: %w", err)
	}

	return res, nil
}

// NewClient returns an operational client that implements the interface
func NewClient() API {
	return &c{}
}

// c implements the API
type c struct{}
