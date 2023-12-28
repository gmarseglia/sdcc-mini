// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: cs-mw.proto

package proto

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

const (
	Front_Choice_FullMethodName = "/mini.Front/Choice"
)

// FrontClient is the client API for Front service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FrontClient interface {
	Choice(ctx context.Context, in *ChoiceBiRequest, opts ...grpc.CallOption) (*ChoiceReply, error)
}

type frontClient struct {
	cc grpc.ClientConnInterface
}

func NewFrontClient(cc grpc.ClientConnInterface) FrontClient {
	return &frontClient{cc}
}

func (c *frontClient) Choice(ctx context.Context, in *ChoiceBiRequest, opts ...grpc.CallOption) (*ChoiceReply, error) {
	out := new(ChoiceReply)
	err := c.cc.Invoke(ctx, Front_Choice_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FrontServer is the server API for Front service.
// All implementations must embed UnimplementedFrontServer
// for forward compatibility
type FrontServer interface {
	Choice(context.Context, *ChoiceBiRequest) (*ChoiceReply, error)
	mustEmbedUnimplementedFrontServer()
}

// UnimplementedFrontServer must be embedded to have forward compatible implementations.
type UnimplementedFrontServer struct {
}

func (UnimplementedFrontServer) Choice(context.Context, *ChoiceBiRequest) (*ChoiceReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Choice not implemented")
}
func (UnimplementedFrontServer) mustEmbedUnimplementedFrontServer() {}

// UnsafeFrontServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FrontServer will
// result in compilation errors.
type UnsafeFrontServer interface {
	mustEmbedUnimplementedFrontServer()
}

func RegisterFrontServer(s grpc.ServiceRegistrar, srv FrontServer) {
	s.RegisterService(&Front_ServiceDesc, srv)
}

func _Front_Choice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChoiceBiRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FrontServer).Choice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Front_Choice_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FrontServer).Choice(ctx, req.(*ChoiceBiRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Front_ServiceDesc is the grpc.ServiceDesc for Front service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Front_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "mini.Front",
	HandlerType: (*FrontServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Choice",
			Handler:    _Front_Choice_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cs-mw.proto",
}

const (
	Back_Choice_FullMethodName = "/mini.Back/Choice"
)

// BackClient is the client API for Back service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BackClient interface {
	Choice(ctx context.Context, in *ChoiceBiRequest, opts ...grpc.CallOption) (*ChoiceReply, error)
}

type backClient struct {
	cc grpc.ClientConnInterface
}

func NewBackClient(cc grpc.ClientConnInterface) BackClient {
	return &backClient{cc}
}

func (c *backClient) Choice(ctx context.Context, in *ChoiceBiRequest, opts ...grpc.CallOption) (*ChoiceReply, error) {
	out := new(ChoiceReply)
	err := c.cc.Invoke(ctx, Back_Choice_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BackServer is the server API for Back service.
// All implementations must embed UnimplementedBackServer
// for forward compatibility
type BackServer interface {
	Choice(context.Context, *ChoiceBiRequest) (*ChoiceReply, error)
	mustEmbedUnimplementedBackServer()
}

// UnimplementedBackServer must be embedded to have forward compatible implementations.
type UnimplementedBackServer struct {
}

func (UnimplementedBackServer) Choice(context.Context, *ChoiceBiRequest) (*ChoiceReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Choice not implemented")
}
func (UnimplementedBackServer) mustEmbedUnimplementedBackServer() {}

// UnsafeBackServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BackServer will
// result in compilation errors.
type UnsafeBackServer interface {
	mustEmbedUnimplementedBackServer()
}

func RegisterBackServer(s grpc.ServiceRegistrar, srv BackServer) {
	s.RegisterService(&Back_ServiceDesc, srv)
}

func _Back_Choice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChoiceBiRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackServer).Choice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Back_Choice_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackServer).Choice(ctx, req.(*ChoiceBiRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Back_ServiceDesc is the grpc.ServiceDesc for Back service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Back_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "mini.Back",
	HandlerType: (*BackServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Choice",
			Handler:    _Back_Choice_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cs-mw.proto",
}

const (
	Master_NotifyActiveWorker_FullMethodName   = "/mini.Master/NotifyActiveWorker"
	Master_NotifyDeactiveWorker_FullMethodName = "/mini.Master/NotifyDeactiveWorker"
)

// MasterClient is the client API for Master service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MasterClient interface {
	NotifyActiveWorker(ctx context.Context, in *NotifyRequest, opts ...grpc.CallOption) (*NotifyReply, error)
	NotifyDeactiveWorker(ctx context.Context, in *NotifyRequest, opts ...grpc.CallOption) (*NotifyReply, error)
}

type masterClient struct {
	cc grpc.ClientConnInterface
}

func NewMasterClient(cc grpc.ClientConnInterface) MasterClient {
	return &masterClient{cc}
}

func (c *masterClient) NotifyActiveWorker(ctx context.Context, in *NotifyRequest, opts ...grpc.CallOption) (*NotifyReply, error) {
	out := new(NotifyReply)
	err := c.cc.Invoke(ctx, Master_NotifyActiveWorker_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *masterClient) NotifyDeactiveWorker(ctx context.Context, in *NotifyRequest, opts ...grpc.CallOption) (*NotifyReply, error) {
	out := new(NotifyReply)
	err := c.cc.Invoke(ctx, Master_NotifyDeactiveWorker_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MasterServer is the server API for Master service.
// All implementations must embed UnimplementedMasterServer
// for forward compatibility
type MasterServer interface {
	NotifyActiveWorker(context.Context, *NotifyRequest) (*NotifyReply, error)
	NotifyDeactiveWorker(context.Context, *NotifyRequest) (*NotifyReply, error)
	mustEmbedUnimplementedMasterServer()
}

// UnimplementedMasterServer must be embedded to have forward compatible implementations.
type UnimplementedMasterServer struct {
}

func (UnimplementedMasterServer) NotifyActiveWorker(context.Context, *NotifyRequest) (*NotifyReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NotifyActiveWorker not implemented")
}
func (UnimplementedMasterServer) NotifyDeactiveWorker(context.Context, *NotifyRequest) (*NotifyReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NotifyDeactiveWorker not implemented")
}
func (UnimplementedMasterServer) mustEmbedUnimplementedMasterServer() {}

// UnsafeMasterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MasterServer will
// result in compilation errors.
type UnsafeMasterServer interface {
	mustEmbedUnimplementedMasterServer()
}

func RegisterMasterServer(s grpc.ServiceRegistrar, srv MasterServer) {
	s.RegisterService(&Master_ServiceDesc, srv)
}

func _Master_NotifyActiveWorker_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NotifyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MasterServer).NotifyActiveWorker(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Master_NotifyActiveWorker_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MasterServer).NotifyActiveWorker(ctx, req.(*NotifyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Master_NotifyDeactiveWorker_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NotifyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MasterServer).NotifyDeactiveWorker(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Master_NotifyDeactiveWorker_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MasterServer).NotifyDeactiveWorker(ctx, req.(*NotifyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Master_ServiceDesc is the grpc.ServiceDesc for Master service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Master_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "mini.Master",
	HandlerType: (*MasterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "NotifyActiveWorker",
			Handler:    _Master_NotifyActiveWorker_Handler,
		},
		{
			MethodName: "NotifyDeactiveWorker",
			Handler:    _Master_NotifyDeactiveWorker_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cs-mw.proto",
}
