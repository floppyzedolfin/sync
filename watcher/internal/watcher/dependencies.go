package watcher

import (
	"context"

	pb "github.com/floppyzedolfin/sync/replica/replica"
	"google.golang.org/grpc"
)

//go:generate mockgen -source dependencies.go -destination mock/dependencies_mock.go

// server lists the methods called by this package
type server interface {
	// File patches a file
	File(ctx context.Context, in *pb.FileRequest, opts ...grpc.CallOption) (*pb.FileResponse, error)
	// Link creates a link
	Link(ctx context.Context, in *pb.LinkRequest, opts ...grpc.CallOption) (*pb.LinkResponse, error)
	// Directory creates a directory
	Directory(ctx context.Context, in *pb.DirectoryRequest, opts ...grpc.CallOption) (*pb.DirectoryResponse, error)
	// Delete an entity on the file system
	Delete(ctx context.Context, in *pb.DeleteRequest, opts ...grpc.CallOption) (*pb.DeleteResponse, error)
}
