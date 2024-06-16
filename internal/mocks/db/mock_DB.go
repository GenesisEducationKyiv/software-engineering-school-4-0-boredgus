// Code generated by mockery v2.42.1. DO NOT EDIT.

package db_mocks

import (
	db "subscription-api/internal/db"

	mock "github.com/stretchr/testify/mock"
)

// DB is an autogenerated mock type for the DB type
type DB struct {
	mock.Mock
}

// DB provides a mock function with given fields:
func (_m *DB) DB() db.Database {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for DB")
	}

	var r0 db.Database
	if rf, ok := ret.Get(0).(func() db.Database); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(db.Database)
		}
	}

	return r0
}

// IsError provides a mock function with given fields: err, errCode
func (_m *DB) IsError(err error, errCode db.Error) bool {
	ret := _m.Called(err, errCode)

	if len(ret) == 0 {
		panic("no return value specified for IsError")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(error, db.Error) bool); ok {
		r0 = rf(err, errCode)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// NewDB creates a new instance of DB. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDB(t interface {
	mock.TestingT
	Cleanup(func())
}) *DB {
	mock := &DB{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}