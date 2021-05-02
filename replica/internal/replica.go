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

// File updates the contents of a file
func (s *server) File(_ context.Context, request *pb.FileRequest) (*pb.FileResponse, error) {
	fmt.Printf("Patching file %s\n", path.Join(s.localReplicaPath, request.FullPath))
	err := ioutil.WriteFile(path.Join(s.localReplicaPath, request.FullPath), []byte(request.FullContents), defaultRights)
	if err != nil {
		return nil, fmt.Errorf("unable to  file %s: %w", request.FullPath, err)
	}
	return &pb.FileResponse{}, nil
}

// Link creates a link to a target
func (s *server) Link(_ context.Context, request *pb.LinkRequest) (*pb.LinkResponse, error) {
	linkLocation := path.Join(s.localReplicaPath, request.FullPath)
	fmt.Printf("Creating a link at %s to %s\n", linkLocation, request.Target)
	err := os.Symlink(request.Target, linkLocation)
	if err != nil {
		return nil, fmt.Errorf("unable to create link %s: %w", request.FullPath, err)
	}
	return &pb.LinkResponse{}, nil
}

// Directory creates a directory
func (s *server) Directory(_ context.Context, request *pb.DirectoryRequest) (*pb.DirectoryResponse, error) {
	fmt.Printf("creating dir %s\n", path.Join(s.localReplicaPath, request.FullPath))
	err := os.Mkdir(path.Join(s.localReplicaPath, request.FullPath), defaultRights)
	if err != nil {
		return nil, fmt.Errorf("unable to create dir %s: %w", request.FullPath, err)
	}
	return &pb.DirectoryResponse{}, nil
}

// Delete removes a file system entity
func (s *server) Delete(_ context.Context, request *pb.DeleteRequest) (*pb.DeleteResponse, error) {
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
