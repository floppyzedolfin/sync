package internal

import (
	"context"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"

	pb "github.com/floppyzedolfin/sync/replica/replica"
)

// NewServer returns an instance of the server that will perform replication to a specified dir
func NewServer(localReplica string) pb.ReplicaServer {
	return &server{
		localReplicaPath: localReplica,
	}
}

const (
	defaultRights = fs.FileMode(0755)
)

// PatchFile updates the contents of a file
func (s *server) PatchFile(ctx context.Context, request *pb.PatchFileRequest) (*pb.PatchFileResponse, error) {
	fmt.Printf("patching file %s\n", path.Join(s.localReplicaPath, request.FullPath))
	err := ioutil.WriteFile(path.Join(s.localReplicaPath, request.FullPath), []byte(request.FullContents), defaultRights)
	if err != nil {
		return nil, fmt.Errorf("unable to patch file %s: %w", request.FullPath, err)
	}
	return &pb.PatchFileResponse{}, nil
}

// CreateDir creates a directory
func (s *server) CreateDir(ctx context.Context, request *pb.CreateDirRequest) (*pb.CreateDirResponse, error) {
	fmt.Printf("creating dir %s\n", path.Join(s.localReplicaPath, request.FullPath))
	err := os.Mkdir(path.Join(s.localReplicaPath, request.FullPath), defaultRights)
	if err != nil {
		return nil, fmt.Errorf("unable to create dir %s: %w", request.FullPath, err)
	}
	return &pb.CreateDirResponse{}, nil
}

// Delete removes a file system entity
func (s *server) Delete(ctx context.Context, request *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	fmt.Printf("deleting %s\n", path.Join(s.localReplicaPath, request.FullPath))
	err := os.Remove(path.Join(s.localReplicaPath, request.FullPath))
	if err != nil {
		return nil, fmt.Errorf("unable to delete at %s: %w", request.FullPath, err)
	}
	return &pb.DeleteResponse{}, nil
}

// server allows for the implementation of the interface
type server struct {
	pb.UnimplementedReplicaServer
	localReplicaPath string
}
