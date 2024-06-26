// Code generated by mockery v2.42.1. DO NOT EDIT.

package client_mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	services "subscription-api/internal/services"
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

// GetAllDispatches provides a mock function with given fields: ctx
func (_m *DispatchServiceClient) GetAllDispatches(ctx context.Context) ([]services.DispatchData, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAllDispatches")
	}

	var r0 []services.DispatchData
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]services.DispatchData, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []services.DispatchData); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]services.DispatchData)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
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
func (_e *DispatchServiceClient_Expecter) GetAllDispatches(ctx interface{}) *DispatchServiceClient_GetAllDispatches_Call {
	return &DispatchServiceClient_GetAllDispatches_Call{Call: _e.mock.On("GetAllDispatches", ctx)}
}

func (_c *DispatchServiceClient_GetAllDispatches_Call) Run(run func(ctx context.Context)) *DispatchServiceClient_GetAllDispatches_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *DispatchServiceClient_GetAllDispatches_Call) Return(_a0 []services.DispatchData, _a1 error) *DispatchServiceClient_GetAllDispatches_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *DispatchServiceClient_GetAllDispatches_Call) RunAndReturn(run func(context.Context) ([]services.DispatchData, error)) *DispatchServiceClient_GetAllDispatches_Call {
	_c.Call.Return(run)
	return _c
}

// SendDispatch provides a mock function with given fields: ctx, dispatchId
func (_m *DispatchServiceClient) SendDispatch(ctx context.Context, dispatchId string) error {
	ret := _m.Called(ctx, dispatchId)

	if len(ret) == 0 {
		panic("no return value specified for SendDispatch")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, dispatchId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DispatchServiceClient_SendDispatch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendDispatch'
type DispatchServiceClient_SendDispatch_Call struct {
	*mock.Call
}

// SendDispatch is a helper method to define mock.On call
//   - ctx context.Context
//   - dispatchId string
func (_e *DispatchServiceClient_Expecter) SendDispatch(ctx interface{}, dispatchId interface{}) *DispatchServiceClient_SendDispatch_Call {
	return &DispatchServiceClient_SendDispatch_Call{Call: _e.mock.On("SendDispatch", ctx, dispatchId)}
}

func (_c *DispatchServiceClient_SendDispatch_Call) Run(run func(ctx context.Context, dispatchId string)) *DispatchServiceClient_SendDispatch_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *DispatchServiceClient_SendDispatch_Call) Return(_a0 error) *DispatchServiceClient_SendDispatch_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DispatchServiceClient_SendDispatch_Call) RunAndReturn(run func(context.Context, string) error) *DispatchServiceClient_SendDispatch_Call {
	_c.Call.Return(run)
	return _c
}

// SubscribeForDispatch provides a mock function with given fields: ctx, email, dispatchId
func (_m *DispatchServiceClient) SubscribeForDispatch(ctx context.Context, email string, dispatchId string) error {
	ret := _m.Called(ctx, email, dispatchId)

	if len(ret) == 0 {
		panic("no return value specified for SubscribeForDispatch")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, email, dispatchId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DispatchServiceClient_SubscribeForDispatch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SubscribeForDispatch'
type DispatchServiceClient_SubscribeForDispatch_Call struct {
	*mock.Call
}

// SubscribeForDispatch is a helper method to define mock.On call
//   - ctx context.Context
//   - email string
//   - dispatchId string
func (_e *DispatchServiceClient_Expecter) SubscribeForDispatch(ctx interface{}, email interface{}, dispatchId interface{}) *DispatchServiceClient_SubscribeForDispatch_Call {
	return &DispatchServiceClient_SubscribeForDispatch_Call{Call: _e.mock.On("SubscribeForDispatch", ctx, email, dispatchId)}
}

func (_c *DispatchServiceClient_SubscribeForDispatch_Call) Run(run func(ctx context.Context, email string, dispatchId string)) *DispatchServiceClient_SubscribeForDispatch_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *DispatchServiceClient_SubscribeForDispatch_Call) Return(_a0 error) *DispatchServiceClient_SubscribeForDispatch_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DispatchServiceClient_SubscribeForDispatch_Call) RunAndReturn(run func(context.Context, string, string) error) *DispatchServiceClient_SubscribeForDispatch_Call {
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
