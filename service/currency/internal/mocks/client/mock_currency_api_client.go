// Code generated by mockery v2.42.1. DO NOT EDIT.

package client_mock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// CurrencyAPIClient is an autogenerated mock type for the CurrencyAPIClient type
type CurrencyAPIClient struct {
	mock.Mock
}

type CurrencyAPIClient_Expecter struct {
	mock *mock.Mock
}

func (_m *CurrencyAPIClient) EXPECT() *CurrencyAPIClient_Expecter {
	return &CurrencyAPIClient_Expecter{mock: &_m.Mock}
}

// Convert provides a mock function with given fields: ctx, baseCcy, targetCcies
func (_m *CurrencyAPIClient) Convert(ctx context.Context, baseCcy string, targetCcies []string) (map[string]float64, error) {
	ret := _m.Called(ctx, baseCcy, targetCcies)

	if len(ret) == 0 {
		panic("no return value specified for Convert")
	}

	var r0 map[string]float64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, []string) (map[string]float64, error)); ok {
		return rf(ctx, baseCcy, targetCcies)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, []string) map[string]float64); ok {
		r0 = rf(ctx, baseCcy, targetCcies)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]float64)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, []string) error); ok {
		r1 = rf(ctx, baseCcy, targetCcies)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CurrencyAPIClient_Convert_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Convert'
type CurrencyAPIClient_Convert_Call struct {
	*mock.Call
}

// Convert is a helper method to define mock.On call
//   - ctx context.Context
//   - baseCcy string
//   - targetCcies []string
func (_e *CurrencyAPIClient_Expecter) Convert(ctx interface{}, baseCcy interface{}, targetCcies interface{}) *CurrencyAPIClient_Convert_Call {
	return &CurrencyAPIClient_Convert_Call{Call: _e.mock.On("Convert", ctx, baseCcy, targetCcies)}
}

func (_c *CurrencyAPIClient_Convert_Call) Run(run func(ctx context.Context, baseCcy string, targetCcies []string)) *CurrencyAPIClient_Convert_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].([]string))
	})
	return _c
}

func (_c *CurrencyAPIClient_Convert_Call) Return(_a0 map[string]float64, _a1 error) *CurrencyAPIClient_Convert_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CurrencyAPIClient_Convert_Call) RunAndReturn(run func(context.Context, string, []string) (map[string]float64, error)) *CurrencyAPIClient_Convert_Call {
	_c.Call.Return(run)
	return _c
}

// NewCurrencyAPIClient creates a new instance of CurrencyAPIClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCurrencyAPIClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *CurrencyAPIClient {
	mock := &CurrencyAPIClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
