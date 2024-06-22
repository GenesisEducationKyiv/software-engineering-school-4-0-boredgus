// Code generated by mockery v2.42.1. DO NOT EDIT.

package services_mocks

import (
	context "context"
	services "subscription-api/internal/services"

	mock "github.com/stretchr/testify/mock"
)

// DispatchService is an autogenerated mock type for the DispatchService type
type DispatchService struct {
	mock.Mock
}

type DispatchService_Expecter struct {
	mock *mock.Mock
}

func (_m *DispatchService) EXPECT() *DispatchService_Expecter {
	return &DispatchService_Expecter{mock: &_m.Mock}
}

// GetAllDispatches provides a mock function with given fields: ctx
func (_m *DispatchService) GetAllDispatches(ctx context.Context) ([]services.DispatchData, error) {
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

// DispatchService_GetAllDispatches_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAllDispatches'
type DispatchService_GetAllDispatches_Call struct {
	*mock.Call
}

// GetAllDispatches is a helper method to define mock.On call
//   - ctx context.Context
func (_e *DispatchService_Expecter) GetAllDispatches(ctx interface{}) *DispatchService_GetAllDispatches_Call {
	return &DispatchService_GetAllDispatches_Call{Call: _e.mock.On("GetAllDispatches", ctx)}
}

func (_c *DispatchService_GetAllDispatches_Call) Run(run func(ctx context.Context)) *DispatchService_GetAllDispatches_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *DispatchService_GetAllDispatches_Call) Return(_a0 []services.DispatchData, _a1 error) *DispatchService_GetAllDispatches_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *DispatchService_GetAllDispatches_Call) RunAndReturn(run func(context.Context) ([]services.DispatchData, error)) *DispatchService_GetAllDispatches_Call {
	_c.Call.Return(run)
	return _c
}

// SendDispatch provides a mock function with given fields: ctx, dispatch
func (_m *DispatchService) SendDispatch(ctx context.Context, dispatch string) error {
	ret := _m.Called(ctx, dispatch)

	if len(ret) == 0 {
		panic("no return value specified for SendDispatch")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, dispatch)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DispatchService_SendDispatch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendDispatch'
type DispatchService_SendDispatch_Call struct {
	*mock.Call
}

// SendDispatch is a helper method to define mock.On call
//   - ctx context.Context
//   - dispatch string
func (_e *DispatchService_Expecter) SendDispatch(ctx interface{}, dispatch interface{}) *DispatchService_SendDispatch_Call {
	return &DispatchService_SendDispatch_Call{Call: _e.mock.On("SendDispatch", ctx, dispatch)}
}

func (_c *DispatchService_SendDispatch_Call) Run(run func(ctx context.Context, dispatch string)) *DispatchService_SendDispatch_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *DispatchService_SendDispatch_Call) Return(_a0 error) *DispatchService_SendDispatch_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DispatchService_SendDispatch_Call) RunAndReturn(run func(context.Context, string) error) *DispatchService_SendDispatch_Call {
	_c.Call.Return(run)
	return _c
}

// SubscribeForDispatch provides a mock function with given fields: ctx, email, dispatch
func (_m *DispatchService) SubscribeForDispatch(ctx context.Context, email string, dispatch string) error {
	ret := _m.Called(ctx, email, dispatch)

	if len(ret) == 0 {
		panic("no return value specified for SubscribeForDispatch")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, email, dispatch)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DispatchService_SubscribeForDispatch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SubscribeForDispatch'
type DispatchService_SubscribeForDispatch_Call struct {
	*mock.Call
}

// SubscribeForDispatch is a helper method to define mock.On call
//   - ctx context.Context
//   - email string
//   - dispatch string
func (_e *DispatchService_Expecter) SubscribeForDispatch(ctx interface{}, email interface{}, dispatch interface{}) *DispatchService_SubscribeForDispatch_Call {
	return &DispatchService_SubscribeForDispatch_Call{Call: _e.mock.On("SubscribeForDispatch", ctx, email, dispatch)}
}

func (_c *DispatchService_SubscribeForDispatch_Call) Run(run func(ctx context.Context, email string, dispatch string)) *DispatchService_SubscribeForDispatch_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *DispatchService_SubscribeForDispatch_Call) Return(_a0 error) *DispatchService_SubscribeForDispatch_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DispatchService_SubscribeForDispatch_Call) RunAndReturn(run func(context.Context, string, string) error) *DispatchService_SubscribeForDispatch_Call {
	_c.Call.Return(run)
	return _c
}

// NewDispatchService creates a new instance of DispatchService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDispatchService(t interface {
	mock.TestingT
	Cleanup(func())
}) *DispatchService {
	mock := &DispatchService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}