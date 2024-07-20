// Code generated by mockery v2.42.1. DO NOT EDIT.

package repo_mock

import (
	context "context"

	service "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/service"
	mock "github.com/stretchr/testify/mock"
)

// SubRepo is an autogenerated mock type for the SubRepo type
type SubRepo struct {
	mock.Mock
}

type SubRepo_Expecter struct {
	mock *mock.Mock
}

func (_m *SubRepo) EXPECT() *SubRepo_Expecter {
	return &SubRepo_Expecter{mock: &_m.Mock}
}

// CreateSubscription provides a mock function with given fields: ctx, args
func (_m *SubRepo) CreateSubscription(ctx context.Context, args service.SubscriptionData) error {
	ret := _m.Called(ctx, args)

	if len(ret) == 0 {
		panic("no return value specified for CreateSubscription")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, service.SubscriptionData) error); ok {
		r0 = rf(ctx, args)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SubRepo_CreateSubscription_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateSubscription'
type SubRepo_CreateSubscription_Call struct {
	*mock.Call
}

// CreateSubscription is a helper method to define mock.On call
//   - ctx context.Context
//   - args service.SubscriptionData
func (_e *SubRepo_Expecter) CreateSubscription(ctx interface{}, args interface{}) *SubRepo_CreateSubscription_Call {
	return &SubRepo_CreateSubscription_Call{Call: _e.mock.On("CreateSubscription", ctx, args)}
}

func (_c *SubRepo_CreateSubscription_Call) Run(run func(ctx context.Context, args service.SubscriptionData)) *SubRepo_CreateSubscription_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(service.SubscriptionData))
	})
	return _c
}

func (_c *SubRepo_CreateSubscription_Call) Return(_a0 error) *SubRepo_CreateSubscription_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *SubRepo_CreateSubscription_Call) RunAndReturn(run func(context.Context, service.SubscriptionData) error) *SubRepo_CreateSubscription_Call {
	_c.Call.Return(run)
	return _c
}

// GetStatusOfSubscription provides a mock function with given fields: ctx, args
func (_m *SubRepo) GetStatusOfSubscription(ctx context.Context, args service.SubscriptionData) (service.SubscriptionStatus, error) {
	ret := _m.Called(ctx, args)

	if len(ret) == 0 {
		panic("no return value specified for GetStatusOfSubscription")
	}

	var r0 service.SubscriptionStatus
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, service.SubscriptionData) (service.SubscriptionStatus, error)); ok {
		return rf(ctx, args)
	}
	if rf, ok := ret.Get(0).(func(context.Context, service.SubscriptionData) service.SubscriptionStatus); ok {
		r0 = rf(ctx, args)
	} else {
		r0 = ret.Get(0).(service.SubscriptionStatus)
	}

	if rf, ok := ret.Get(1).(func(context.Context, service.SubscriptionData) error); ok {
		r1 = rf(ctx, args)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SubRepo_GetStatusOfSubscription_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetStatusOfSubscription'
type SubRepo_GetStatusOfSubscription_Call struct {
	*mock.Call
}

// GetStatusOfSubscription is a helper method to define mock.On call
//   - ctx context.Context
//   - args service.SubscriptionData
func (_e *SubRepo_Expecter) GetStatusOfSubscription(ctx interface{}, args interface{}) *SubRepo_GetStatusOfSubscription_Call {
	return &SubRepo_GetStatusOfSubscription_Call{Call: _e.mock.On("GetStatusOfSubscription", ctx, args)}
}

func (_c *SubRepo_GetStatusOfSubscription_Call) Run(run func(ctx context.Context, args service.SubscriptionData)) *SubRepo_GetStatusOfSubscription_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(service.SubscriptionData))
	})
	return _c
}

func (_c *SubRepo_GetStatusOfSubscription_Call) Return(_a0 service.SubscriptionStatus, _a1 error) *SubRepo_GetStatusOfSubscription_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SubRepo_GetStatusOfSubscription_Call) RunAndReturn(run func(context.Context, service.SubscriptionData) (service.SubscriptionStatus, error)) *SubRepo_GetStatusOfSubscription_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateSubscriptionStatus provides a mock function with given fields: ctx, args, status
func (_m *SubRepo) UpdateSubscriptionStatus(ctx context.Context, args service.SubscriptionData, status service.SubscriptionStatus) error {
	ret := _m.Called(ctx, args, status)

	if len(ret) == 0 {
		panic("no return value specified for UpdateSubscriptionStatus")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, service.SubscriptionData, service.SubscriptionStatus) error); ok {
		r0 = rf(ctx, args, status)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SubRepo_UpdateSubscriptionStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateSubscriptionStatus'
type SubRepo_UpdateSubscriptionStatus_Call struct {
	*mock.Call
}

// UpdateSubscriptionStatus is a helper method to define mock.On call
//   - ctx context.Context
//   - args service.SubscriptionData
//   - status service.SubscriptionStatus
func (_e *SubRepo_Expecter) UpdateSubscriptionStatus(ctx interface{}, args interface{}, status interface{}) *SubRepo_UpdateSubscriptionStatus_Call {
	return &SubRepo_UpdateSubscriptionStatus_Call{Call: _e.mock.On("UpdateSubscriptionStatus", ctx, args, status)}
}

func (_c *SubRepo_UpdateSubscriptionStatus_Call) Run(run func(ctx context.Context, args service.SubscriptionData, status service.SubscriptionStatus)) *SubRepo_UpdateSubscriptionStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(service.SubscriptionData), args[2].(service.SubscriptionStatus))
	})
	return _c
}

func (_c *SubRepo_UpdateSubscriptionStatus_Call) Return(_a0 error) *SubRepo_UpdateSubscriptionStatus_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *SubRepo_UpdateSubscriptionStatus_Call) RunAndReturn(run func(context.Context, service.SubscriptionData, service.SubscriptionStatus) error) *SubRepo_UpdateSubscriptionStatus_Call {
	_c.Call.Return(run)
	return _c
}

// NewSubRepo creates a new instance of SubRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSubRepo(t interface {
	mock.TestingT
	Cleanup(func())
}) *SubRepo {
	mock := &SubRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
