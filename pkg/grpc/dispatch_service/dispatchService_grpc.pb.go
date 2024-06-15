// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.26.1
// source: pkg/grpc/dispatch_service/dispatchService.proto

package __

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

// DispatchServiceClient is the client API for DispatchService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DispatchServiceClient interface {
	SubscribeForDispatch(ctx context.Context, in *SubscribeForDispatchRequest, opts ...grpc.CallOption) (*SubscribeForDispatchResponse, error)
	SendDispatch(ctx context.Context, in *SendDispatchRequest, opts ...grpc.CallOption) (*SendDispatchResponse, error)
	GetAllDispatches(ctx context.Context, in *GetAllDispatchesRequest, opts ...grpc.CallOption) (*GetAllDispatchesResponse, error)
}

type dispatchServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDispatchServiceClient(cc grpc.ClientConnInterface) DispatchServiceClient {
	return &dispatchServiceClient{cc}
}

func (c *dispatchServiceClient) SubscribeForDispatch(ctx context.Context, in *SubscribeForDispatchRequest, opts ...grpc.CallOption) (*SubscribeForDispatchResponse, error) {
	out := new(SubscribeForDispatchResponse)
	err := c.cc.Invoke(ctx, "/main.DispatchService/SubscribeForDispatch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dispatchServiceClient) SendDispatch(ctx context.Context, in *SendDispatchRequest, opts ...grpc.CallOption) (*SendDispatchResponse, error) {
	out := new(SendDispatchResponse)
	err := c.cc.Invoke(ctx, "/main.DispatchService/SendDispatch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dispatchServiceClient) GetAllDispatches(ctx context.Context, in *GetAllDispatchesRequest, opts ...grpc.CallOption) (*GetAllDispatchesResponse, error) {
	out := new(GetAllDispatchesResponse)
	err := c.cc.Invoke(ctx, "/main.DispatchService/GetAllDispatches", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DispatchServiceServer is the server API for DispatchService service.
// All implementations must embed UnimplementedDispatchServiceServer
// for forward compatibility
type DispatchServiceServer interface {
	SubscribeForDispatch(context.Context, *SubscribeForDispatchRequest) (*SubscribeForDispatchResponse, error)
	SendDispatch(context.Context, *SendDispatchRequest) (*SendDispatchResponse, error)
	GetAllDispatches(context.Context, *GetAllDispatchesRequest) (*GetAllDispatchesResponse, error)
	mustEmbedUnimplementedDispatchServiceServer()
}

// UnimplementedDispatchServiceServer must be embedded to have forward compatible implementations.
type UnimplementedDispatchServiceServer struct {
}

func (UnimplementedDispatchServiceServer) SubscribeForDispatch(context.Context, *SubscribeForDispatchRequest) (*SubscribeForDispatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubscribeForDispatch not implemented")
}
func (UnimplementedDispatchServiceServer) SendDispatch(context.Context, *SendDispatchRequest) (*SendDispatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendDispatch not implemented")
}
func (UnimplementedDispatchServiceServer) GetAllDispatches(context.Context, *GetAllDispatchesRequest) (*GetAllDispatchesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllDispatches not implemented")
}
func (UnimplementedDispatchServiceServer) mustEmbedUnimplementedDispatchServiceServer() {}

// UnsafeDispatchServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DispatchServiceServer will
// result in compilation errors.
type UnsafeDispatchServiceServer interface {
	mustEmbedUnimplementedDispatchServiceServer()
}

func RegisterDispatchServiceServer(s grpc.ServiceRegistrar, srv DispatchServiceServer) {
	s.RegisterService(&DispatchService_ServiceDesc, srv)
}

func _DispatchService_SubscribeForDispatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubscribeForDispatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DispatchServiceServer).SubscribeForDispatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.DispatchService/SubscribeForDispatch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DispatchServiceServer).SubscribeForDispatch(ctx, req.(*SubscribeForDispatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DispatchService_SendDispatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendDispatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DispatchServiceServer).SendDispatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.DispatchService/SendDispatch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DispatchServiceServer).SendDispatch(ctx, req.(*SendDispatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DispatchService_GetAllDispatches_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllDispatchesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DispatchServiceServer).GetAllDispatches(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.DispatchService/GetAllDispatches",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DispatchServiceServer).GetAllDispatches(ctx, req.(*GetAllDispatchesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DispatchService_ServiceDesc is the grpc.ServiceDesc for DispatchService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DispatchService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "main.DispatchService",
	HandlerType: (*DispatchServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SubscribeForDispatch",
			Handler:    _DispatchService_SubscribeForDispatch_Handler,
		},
		{
			MethodName: "SendDispatch",
			Handler:    _DispatchService_SendDispatch_Handler,
		},
		{
			MethodName: "GetAllDispatches",
			Handler:    _DispatchService_GetAllDispatches_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/grpc/dispatch_service/dispatchService.proto",
}
