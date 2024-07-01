// Code generated by mockery v2.42.1. DO NOT EDIT.

package repo_mock

import (
	context "context"

	repo "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/repo"
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
func (_m *SubRepo) CreateSubscription(ctx context.Context, args repo.SubscriptionData) error {
	ret := _m.Called(ctx, args)

	if len(ret) == 0 {
		panic("no return value specified for CreateSubscription")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, repo.SubscriptionData) error); ok {
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
//   - args repo.SubscriptionData
func (_e *SubRepo_Expecter) CreateSubscription(ctx interface{}, args interface{}) *SubRepo_CreateSubscription_Call {
	return &SubRepo_CreateSubscription_Call{Call: _e.mock.On("CreateSubscription", ctx, args)}
}

func (_c *SubRepo_CreateSubscription_Call) Run(run func(ctx context.Context, args repo.SubscriptionData)) *SubRepo_CreateSubscription_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(repo.SubscriptionData))
	})
	return _c
}

func (_c *SubRepo_CreateSubscription_Call) Return(_a0 error) *SubRepo_CreateSubscription_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *SubRepo_CreateSubscription_Call) RunAndReturn(run func(context.Context, repo.SubscriptionData) error) *SubRepo_CreateSubscription_Call {
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