// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package replica

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ReplicaClient is the client API for Replica service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ReplicaClient interface {
	// File patches a file
	File(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (*FileResponse, error)
	// Directory creates a directory
	Directory(ctx context.Context, in *DirectoryRequest, opts ...grpc.CallOption) (*DirectoryResponse, error)
	// Link creates a link
	Link(ctx context.Context, in *LinkRequest, opts ...grpc.CallOption) (*LinkResponse, error)
	// Delete an entity on the file system
	Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error)
}

type replicaClient struct {
	cc grpc.ClientConnInterface
}

func NewReplicaClient(cc grpc.ClientConnInterface) ReplicaClient {
	return &replicaClient{cc}
}

func (c *replicaClient) File(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (*FileResponse, error) {
	out := new(FileResponse)
	err := c.cc.Invoke(ctx, "/replica.Replica/File", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *replicaClient) Directory(ctx context.Context, in *DirectoryRequest, opts ...grpc.CallOption) (*DirectoryResponse, error) {
	out := new(DirectoryResponse)
	err := c.cc.Invoke(ctx, "/replica.Replica/Directory", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *replicaClient) Link(ctx context.Context, in *LinkRequest, opts ...grpc.CallOption) (*LinkResponse, error) {
	out := new(LinkResponse)
	err := c.cc.Invoke(ctx, "/replica.Replica/Link", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *replicaClient) Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error) {
	out := new(DeleteResponse)
	err := c.cc.Invoke(ctx, "/replica.Replica/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ReplicaServer is the server API for Replica service.
// All implementations must embed UnimplementedReplicaServer
// for forward compatibility
type ReplicaServer interface {
	// File patches a file
	File(context.Context, *FileRequest) (*FileResponse, error)
	// Directory creates a directory
	Directory(context.Context, *DirectoryRequest) (*DirectoryResponse, error)
	// Link creates a link
	Link(context.Context, *LinkRequest) (*LinkResponse, error)
	// Delete an entity on the file system
	Delete(context.Context, *DeleteRequest) (*DeleteResponse, error)
	mustEmbedUnimplementedReplicaServer()
}

// UnimplementedReplicaServer must be embedded to have forward compatible implementations.
type UnimplementedReplicaServer struct {
}

func (UnimplementedReplicaServer) File(context.Context, *FileRequest) (*FileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method File not implemented")
}
func (UnimplementedReplicaServer) Directory(context.Context, *DirectoryRequest) (*DirectoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Directory not implemented")
}
func (UnimplementedReplicaServer) Link(context.Context, *LinkRequest) (*LinkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Link not implemented")
}
func (UnimplementedReplicaServer) Delete(context.Context, *DeleteRequest) (*DeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedReplicaServer) mustEmbedUnimplementedReplicaServer() {}

// UnsafeReplicaServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ReplicaServer will
// result in compilation errors.
type UnsafeReplicaServer interface {
	mustEmbedUnimplementedReplicaServer()
}

func RegisterReplicaServer(s grpc.ServiceRegistrar, srv ReplicaServer) {
	s.RegisterService(&Replica_ServiceDesc, srv)
}

func _Replica_File_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReplicaServer).File(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/replica.Replica/File",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReplicaServer).File(ctx, req.(*FileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Replica_Directory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DirectoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReplicaServer).Directory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/replica.Replica/Directory",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReplicaServer).Directory(ctx, req.(*DirectoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Replica_Link_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LinkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReplicaServer).Link(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/replica.Replica/Link",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReplicaServer).Link(ctx, req.(*LinkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Replica_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReplicaServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/replica.Replica/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReplicaServer).Delete(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Replica_ServiceDesc is the grpc.ServiceDesc for Replica service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Replica_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "replica.Replica",
	HandlerType: (*ReplicaServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "File",
			Handler:    _Replica_File_Handler,
		},
		{
			MethodName: "Directory",
			Handler:    _Replica_Directory_Handler,
		},
		{
			MethodName: "Link",
			Handler:    _Replica_Link_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _Replica_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "replica.proto",
}
