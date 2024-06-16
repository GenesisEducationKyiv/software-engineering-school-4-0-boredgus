// Code generated by mockery v2.42.1. DO NOT EDIT.

package mailing_mocks

import (
	mail "github.com/go-mail/mail"

	mock "github.com/stretchr/testify/mock"
)

// Dialer is an autogenerated mock type for the Dialer type
type Dialer struct {
	mock.Mock
}

type Dialer_Expecter struct {
	mock *mock.Mock
}

func (_m *Dialer) EXPECT() *Dialer_Expecter {
	return &Dialer_Expecter{mock: &_m.Mock}
}

// DialAndSend provides a mock function with given fields: m
func (_m *Dialer) DialAndSend(m ...*mail.Message) error {
	_va := make([]interface{}, len(m))
	for _i := range m {
		_va[_i] = m[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for DialAndSend")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(...*mail.Message) error); ok {
		r0 = rf(m...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Dialer_DialAndSend_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DialAndSend'
type Dialer_DialAndSend_Call struct {
	*mock.Call
}

// DialAndSend is a helper method to define mock.On call
//   - m ...*mail.Message
func (_e *Dialer_Expecter) DialAndSend(m ...interface{}) *Dialer_DialAndSend_Call {
	return &Dialer_DialAndSend_Call{Call: _e.mock.On("DialAndSend",
		append([]interface{}{}, m...)...)}
}

func (_c *Dialer_DialAndSend_Call) Run(run func(m ...*mail.Message)) *Dialer_DialAndSend_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]*mail.Message, len(args)-0)
		for i, a := range args[0:] {
			if a != nil {
				variadicArgs[i] = a.(*mail.Message)
			}
		}
		run(variadicArgs...)
	})
	return _c
}

func (_c *Dialer_DialAndSend_Call) Return(_a0 error) *Dialer_DialAndSend_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Dialer_DialAndSend_Call) RunAndReturn(run func(...*mail.Message) error) *Dialer_DialAndSend_Call {
	_c.Call.Return(run)
	return _c
}

// NewDialer creates a new instance of Dialer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDialer(t interface {
	mock.TestingT
	Cleanup(func())
}) *Dialer {
	mock := &Dialer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
