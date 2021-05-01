package watcher

import (
	"context"

	pb "github.com/floppyzedolfin/sync/replica/replica"
	"google.golang.org/grpc"
)

//go:generate mockgen -source dependencies.go -destination mock/dependencies_mock.go

// server lists the methods called by this package
type server interface {
	// PatchFile patches a file
	PatchFile(ctx context.Context, in *pb.PatchFileRequest, opts ...grpc.CallOption) (*pb.PatchFileResponse, error)
	// CreateDir creates a directory
	CreateDir(ctx context.Context, in *pb.CreateDirRequest, opts ...grpc.CallOption) (*pb.CreateDirResponse, error)
	// Delete an entity on the file system
	Delete(ctx context.Context, in *pb.DeleteRequest, opts ...grpc.CallOption) (*pb.DeleteResponse, error)
}
