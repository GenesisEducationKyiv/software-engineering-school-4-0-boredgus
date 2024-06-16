// Code generated by mockery v2.42.1. DO NOT EDIT.

package client_mocks

import (
	context "context"
	__ "subscription-api/pkg/grpc/dispatch_service"

	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"
)

// DispatchServiceClient is an autogenerated mock type for the DispatchServiceClient type
type DispatchServiceClient struct {
	mock.Mock
}

type DispatchServiceClient_Expecter struct {
	mock *mock.Mock
}

func (_m *DispatchServiceClient) EXPECT() *DispatchServiceClient_Expecter {
	return &DispatchServiceClient_Expecter{mock: &_m.Mock}
}

// GetAllDispatches provides a mock function with given fields: ctx, in, opts
func (_m *DispatchServiceClient) GetAllDispatches(ctx context.Context, in *__.GetAllDispatchesRequest, opts ...grpc.CallOption) (*__.GetAllDispatchesResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetAllDispatches")
	}

	var r0 *__.GetAllDispatchesResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *__.GetAllDispatchesRequest, ...grpc.CallOption) (*__.GetAllDispatchesResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *__.GetAllDispatchesRequest, ...grpc.CallOption) *__.GetAllDispatchesResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*__.GetAllDispatchesResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *__.GetAllDispatchesRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DispatchServiceClient_GetAllDispatches_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAllDispatches'
type DispatchServiceClient_GetAllDispatches_Call struct {
	*mock.Call
}

// GetAllDispatches is a helper method to define mock.On call
//   - ctx context.Context
//   - in *__.GetAllDispatchesRequest
//   - opts ...grpc.CallOption
func (_e *DispatchServiceClient_Expecter) GetAllDispatches(ctx interface{}, in interface{}, opts ...interface{}) *DispatchServiceClient_GetAllDispatches_Call {
	return &DispatchServiceClient_GetAllDispatches_Call{Call: _e.mock.On("GetAllDispatches",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *DispatchServiceClient_GetAllDispatches_Call) Run(run func(ctx context.Context, in *__.GetAllDispatchesRequest, opts ...grpc.CallOption)) *DispatchServiceClient_GetAllDispatches_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*__.GetAllDispatchesRequest), variadicArgs...)
	})
	return _c
}

func (_c *DispatchServiceClient_GetAllDispatches_Call) Return(_a0 *__.GetAllDispatchesResponse, _a1 error) *DispatchServiceClient_GetAllDispatches_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *DispatchServiceClient_GetAllDispatches_Call) RunAndReturn(run func(context.Context, *__.GetAllDispatchesRequest, ...grpc.CallOption) (*__.GetAllDispatchesResponse, error)) *DispatchServiceClient_GetAllDispatches_Call {
	_c.Call.Return(run)
	return _c
}

// SendDispatch provides a mock function with given fields: ctx, in, opts
func (_m *DispatchServiceClient) SendDispatch(ctx context.Context, in *__.SendDispatchRequest, opts ...grpc.CallOption) (*__.SendDispatchResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for SendDispatch")
	}

	var r0 *__.SendDispatchResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *__.SendDispatchRequest, ...grpc.CallOption) (*__.SendDispatchResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *__.SendDispatchRequest, ...grpc.CallOption) *__.SendDispatchResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*__.SendDispatchResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *__.SendDispatchRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DispatchServiceClient_SendDispatch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendDispatch'
type DispatchServiceClient_SendDispatch_Call struct {
	*mock.Call
}

// SendDispatch is a helper method to define mock.On call
//   - ctx context.Context
//   - in *__.SendDispatchRequest
//   - opts ...grpc.CallOption
func (_e *DispatchServiceClient_Expecter) SendDispatch(ctx interface{}, in interface{}, opts ...interface{}) *DispatchServiceClient_SendDispatch_Call {
	return &DispatchServiceClient_SendDispatch_Call{Call: _e.mock.On("SendDispatch",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *DispatchServiceClient_SendDispatch_Call) Run(run func(ctx context.Context, in *__.SendDispatchRequest, opts ...grpc.CallOption)) *DispatchServiceClient_SendDispatch_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*__.SendDispatchRequest), variadicArgs...)
	})
	return _c
}

func (_c *DispatchServiceClient_SendDispatch_Call) Return(_a0 *__.SendDispatchResponse, _a1 error) *DispatchServiceClient_SendDispatch_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *DispatchServiceClient_SendDispatch_Call) RunAndReturn(run func(context.Context, *__.SendDispatchRequest, ...grpc.CallOption) (*__.SendDispatchResponse, error)) *DispatchServiceClient_SendDispatch_Call {
	_c.Call.Return(run)
	return _c
}

// SubscribeForDispatch provides a mock function with given fields: ctx, in, opts
func (_m *DispatchServiceClient) SubscribeForDispatch(ctx context.Context, in *__.SubscribeForDispatchRequest, opts ...grpc.CallOption) (*__.SubscribeForDispatchResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for SubscribeForDispatch")
	}

	var r0 *__.SubscribeForDispatchResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *__.SubscribeForDispatchRequest, ...grpc.CallOption) (*__.SubscribeForDispatchResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *__.SubscribeForDispatchRequest, ...grpc.CallOption) *__.SubscribeForDispatchResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*__.SubscribeForDispatchResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *__.SubscribeForDispatchRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DispatchServiceClient_SubscribeForDispatch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SubscribeForDispatch'
type DispatchServiceClient_SubscribeForDispatch_Call struct {
	*mock.Call
}

// SubscribeForDispatch is a helper method to define mock.On call
//   - ctx context.Context
//   - in *__.SubscribeForDispatchRequest
//   - opts ...grpc.CallOption
func (_e *DispatchServiceClient_Expecter) SubscribeForDispatch(ctx interface{}, in interface{}, opts ...interface{}) *DispatchServiceClient_SubscribeForDispatch_Call {
	return &DispatchServiceClient_SubscribeForDispatch_Call{Call: _e.mock.On("SubscribeForDispatch",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *DispatchServiceClient_SubscribeForDispatch_Call) Run(run func(ctx context.Context, in *__.SubscribeForDispatchRequest, opts ...grpc.CallOption)) *DispatchServiceClient_SubscribeForDispatch_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*__.SubscribeForDispatchRequest), variadicArgs...)
	})
	return _c
}

func (_c *DispatchServiceClient_SubscribeForDispatch_Call) Return(_a0 *__.SubscribeForDispatchResponse, _a1 error) *DispatchServiceClient_SubscribeForDispatch_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *DispatchServiceClient_SubscribeForDispatch_Call) RunAndReturn(run func(context.Context, *__.SubscribeForDispatchRequest, ...grpc.CallOption) (*__.SubscribeForDispatchResponse, error)) *DispatchServiceClient_SubscribeForDispatch_Call {
	_c.Call.Return(run)
	return _c
}

// NewDispatchServiceClient creates a new instance of DispatchServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDispatchServiceClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *DispatchServiceClient {
	mock := &DispatchServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}