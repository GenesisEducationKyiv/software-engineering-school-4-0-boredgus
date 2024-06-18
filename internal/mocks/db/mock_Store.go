// Code generated by mockery v2.42.1. DO NOT EDIT.

package db_mocks

import (
	context "context"
	db "subscription-api/internal/db"

	mock "github.com/stretchr/testify/mock"
)

// Store is an autogenerated mock type for the Store type
type Store struct {
	mock.Mock
}

type Store_Expecter struct {
	mock *mock.Mock
}

func (_m *Store) EXPECT() *Store_Expecter {
	return &Store_Expecter{mock: &_m.Mock}
}

// WithTx provides a mock function with given fields: ctx, f
func (_m *Store) WithTx(ctx context.Context, f func(db.DB) error) error {
	ret := _m.Called(ctx, f)

	if len(ret) == 0 {
		panic("no return value specified for WithTx")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, func(db.DB) error) error); ok {
		r0 = rf(ctx, f)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Store_WithTx_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithTx'
type Store_WithTx_Call struct {
	*mock.Call
}

// WithTx is a helper method to define mock.On call
//   - ctx context.Context
//   - f func(db.DB) error
func (_e *Store_Expecter) WithTx(ctx interface{}, f interface{}) *Store_WithTx_Call {
	return &Store_WithTx_Call{Call: _e.mock.On("WithTx", ctx, f)}
}

func (_c *Store_WithTx_Call) Run(run func(ctx context.Context, f func(db.DB) error)) *Store_WithTx_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(func(db.DB) error))
	})
	return _c
}

func (_c *Store_WithTx_Call) Return(_a0 error) *Store_WithTx_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Store_WithTx_Call) RunAndReturn(run func(context.Context, func(db.DB) error) error) *Store_WithTx_Call {
	_c.Call.Return(run)
	return _c
}

// NewStore creates a new instance of Store. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *Store {
	mock := &Store{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
