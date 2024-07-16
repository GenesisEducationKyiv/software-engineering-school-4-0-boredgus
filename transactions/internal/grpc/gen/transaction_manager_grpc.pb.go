// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.26.1
// source: services/transaction_manager.proto

package __

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// TransactionManagerClient is the client API for TransactionManager service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TransactionManagerClient interface {
	SubscribeForDispatch(ctx context.Context, in *SubscribeForDispatchRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	UnsubscribeFromDispatch(ctx context.Context, in *UnsubscribeFromDispatchRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type transactionManagerClient struct {
	cc grpc.ClientConnInterface
}

func NewTransactionManagerClient(cc grpc.ClientConnInterface) TransactionManagerClient {
	return &transactionManagerClient{cc}
}

func (c *transactionManagerClient) SubscribeForDispatch(ctx context.Context, in *SubscribeForDispatchRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/services.TransactionManager/SubscribeForDispatch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *transactionManagerClient) UnsubscribeFromDispatch(ctx context.Context, in *UnsubscribeFromDispatchRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/services.TransactionManager/UnsubscribeFromDispatch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TransactionManagerServer is the server API for TransactionManager service.
// All implementations must embed UnimplementedTransactionManagerServer
// for forward compatibility
type TransactionManagerServer interface {
	SubscribeForDispatch(context.Context, *SubscribeForDispatchRequest) (*emptypb.Empty, error)
	UnsubscribeFromDispatch(context.Context, *UnsubscribeFromDispatchRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedTransactionManagerServer()
}

// UnimplementedTransactionManagerServer must be embedded to have forward compatible implementations.
type UnimplementedTransactionManagerServer struct {
}

func (UnimplementedTransactionManagerServer) SubscribeForDispatch(context.Context, *SubscribeForDispatchRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubscribeForDispatch not implemented")
}
func (UnimplementedTransactionManagerServer) UnsubscribeFromDispatch(context.Context, *UnsubscribeFromDispatchRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnsubscribeFromDispatch not implemented")
}
func (UnimplementedTransactionManagerServer) mustEmbedUnimplementedTransactionManagerServer() {}

// UnsafeTransactionManagerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TransactionManagerServer will
// result in compilation errors.
type UnsafeTransactionManagerServer interface {
	mustEmbedUnimplementedTransactionManagerServer()
}

func RegisterTransactionManagerServer(s grpc.ServiceRegistrar, srv TransactionManagerServer) {
	s.RegisterService(&TransactionManager_ServiceDesc, srv)
}

func _TransactionManager_SubscribeForDispatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubscribeForDispatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransactionManagerServer).SubscribeForDispatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.TransactionManager/SubscribeForDispatch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransactionManagerServer).SubscribeForDispatch(ctx, req.(*SubscribeForDispatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TransactionManager_UnsubscribeFromDispatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnsubscribeFromDispatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransactionManagerServer).UnsubscribeFromDispatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.TransactionManager/UnsubscribeFromDispatch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransactionManagerServer).UnsubscribeFromDispatch(ctx, req.(*UnsubscribeFromDispatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TransactionManager_ServiceDesc is the grpc.ServiceDesc for TransactionManager service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TransactionManager_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "services.TransactionManager",
	HandlerType: (*TransactionManagerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SubscribeForDispatch",
			Handler:    _TransactionManager_SubscribeForDispatch_Handler,
		},
		{
			MethodName: "UnsubscribeFromDispatch",
			Handler:    _TransactionManager_UnsubscribeFromDispatch_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "services/transaction_manager.proto",
}