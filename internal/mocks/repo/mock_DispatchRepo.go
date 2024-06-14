// Code generated by mockery v2.42.1. DO NOT EDIT.

package repo_mocks

import (
	context "context"
	db "subscription-api/internal/db"

	mock "github.com/stretchr/testify/mock"
)

// DispatchRepo is an autogenerated mock type for the DispatchRepo type
type DispatchRepo struct {
	mock.Mock
}

type DispatchRepo_Expecter struct {
	mock *mock.Mock
}

func (_m *DispatchRepo) EXPECT() *DispatchRepo_Expecter {
	return &DispatchRepo_Expecter{mock: &_m.Mock}
}

// GetByID provides a mock function with given fields: ctx, _a1, dispatchId
func (_m *DispatchRepo) GetByID(ctx context.Context, _a1 db.DB, dispatchId string) (db.DispatchData, error) {
	ret := _m.Called(ctx, _a1, dispatchId)

	if len(ret) == 0 {
		panic("no return value specified for GetByID")
	}

	var r0 db.DispatchData
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, db.DB, string) (db.DispatchData, error)); ok {
		return rf(ctx, _a1, dispatchId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, db.DB, string) db.DispatchData); ok {
		r0 = rf(ctx, _a1, dispatchId)
	} else {
		r0 = ret.Get(0).(db.DispatchData)
	}

	if rf, ok := ret.Get(1).(func(context.Context, db.DB, string) error); ok {
		r1 = rf(ctx, _a1, dispatchId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DispatchRepo_GetByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByID'
type DispatchRepo_GetByID_Call struct {
	*mock.Call
}

// GetByID is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 db.DB
//   - dispatchId string
func (_e *DispatchRepo_Expecter) GetByID(ctx interface{}, _a1 interface{}, dispatchId interface{}) *DispatchRepo_GetByID_Call {
	return &DispatchRepo_GetByID_Call{Call: _e.mock.On("GetByID", ctx, _a1, dispatchId)}
}

func (_c *DispatchRepo_GetByID_Call) Run(run func(ctx context.Context, _a1 db.DB, dispatchId string)) *DispatchRepo_GetByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(db.DB), args[2].(string))
	})
	return _c
}

func (_c *DispatchRepo_GetByID_Call) Return(_a0 db.DispatchData, _a1 error) *DispatchRepo_GetByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *DispatchRepo_GetByID_Call) RunAndReturn(run func(context.Context, db.DB, string) (db.DispatchData, error)) *DispatchRepo_GetByID_Call {
	_c.Call.Return(run)
	return _c
}

// NewDispatchRepo creates a new instance of DispatchRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDispatchRepo(t interface {
	mock.TestingT
	Cleanup(func())
}) *DispatchRepo {
	mock := &DispatchRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}