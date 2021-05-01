package watcher

import (
	"context"

	pb "github.com/floppyzedolfin/sync/replica/replica"
)

type server interface {
	PatchFile(context.Context, pb.PatchFileRequest) (*pb.PatchFileResponse, error)
	CreateDir(context.Context, pb.CreateDirRequest) (*pb.CreateDirResponse, error)
	Delete(context.Context, pb.DeleteRequest) (*pb.DeleteResponse, error)
}
