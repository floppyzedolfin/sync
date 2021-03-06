// Code generated by MockGen. DO NOT EDIT.
// Source: dependencies.go

// Package mock_watcher is a generated GoMock package.
package mock_watcher

import (
	context "context"
	reflect "reflect"

	replica "github.com/floppyzedolfin/sync/replica/replica"
	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
)

// Mockserver is a mock of server interface.
type Mockserver struct {
	ctrl     *gomock.Controller
	recorder *MockserverMockRecorder
}

// MockserverMockRecorder is the mock recorder for Mockserver.
type MockserverMockRecorder struct {
	mock *Mockserver
}

// NewMockserver creates a new mock instance.
func NewMockserver(ctrl *gomock.Controller) *Mockserver {
	mock := &Mockserver{ctrl: ctrl}
	mock.recorder = &MockserverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockserver) EXPECT() *MockserverMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *Mockserver) Delete(ctx context.Context, in *replica.DeleteRequest, opts ...grpc.CallOption) (*replica.DeleteResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Delete", varargs...)
	ret0, _ := ret[0].(*replica.DeleteResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockserverMockRecorder) Delete(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*Mockserver)(nil).Delete), varargs...)
}

// Directory mocks base method.
func (m *Mockserver) Directory(ctx context.Context, in *replica.DirectoryRequest, opts ...grpc.CallOption) (*replica.DirectoryResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Directory", varargs...)
	ret0, _ := ret[0].(*replica.DirectoryResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Directory indicates an expected call of Directory.
func (mr *MockserverMockRecorder) Directory(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Directory", reflect.TypeOf((*Mockserver)(nil).Directory), varargs...)
}

// File mocks base method.
func (m *Mockserver) File(ctx context.Context, in *replica.FileRequest, opts ...grpc.CallOption) (*replica.FileResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "File", varargs...)
	ret0, _ := ret[0].(*replica.FileResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// File indicates an expected call of File.
func (mr *MockserverMockRecorder) File(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "File", reflect.TypeOf((*Mockserver)(nil).File), varargs...)
}

// Link mocks base method.
func (m *Mockserver) Link(ctx context.Context, in *replica.LinkRequest, opts ...grpc.CallOption) (*replica.LinkResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Link", varargs...)
	ret0, _ := ret[0].(*replica.LinkResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Link indicates an expected call of Link.
func (mr *MockserverMockRecorder) Link(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Link", reflect.TypeOf((*Mockserver)(nil).Link), varargs...)
}
