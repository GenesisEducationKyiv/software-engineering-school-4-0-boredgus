// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	context "context"
	entities "subscription-api/internal/entities"

	mock "github.com/stretchr/testify/mock"

	services "subscription-api/internal/services/currency"
)

// CurrencyService is an autogenerated mock type for the CurrencyService type
type CurrencyService struct {
	mock.Mock
}

type CurrencyService_Expecter struct {
	mock *mock.Mock
}

func (_m *CurrencyService) EXPECT() *CurrencyService_Expecter {
	return &CurrencyService_Expecter{mock: &_m.Mock}
}

// Convert provides a mock function with given fields: ctx, params
func (_m *CurrencyService) Convert(ctx context.Context, params services.ConvertParams) (map[entities.Currency]float64, error) {
	ret := _m.Called(ctx, params)

	if len(ret) == 0 {
		panic("no return value specified for Convert")
	}

	var r0 map[entities.Currency]float64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, services.ConvertParams) (map[entities.Currency]float64, error)); ok {
		return rf(ctx, params)
	}
	if rf, ok := ret.Get(0).(func(context.Context, services.ConvertParams) map[entities.Currency]float64); ok {
		r0 = rf(ctx, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[entities.Currency]float64)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, services.ConvertParams) error); ok {
		r1 = rf(ctx, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CurrencyService_Convert_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Convert'
type CurrencyService_Convert_Call struct {
	*mock.Call
}

// Convert is a helper method to define mock.On call
//   - ctx context.Context
//   - params services.ConvertParams
func (_e *CurrencyService_Expecter) Convert(ctx interface{}, params interface{}) *CurrencyService_Convert_Call {
	return &CurrencyService_Convert_Call{Call: _e.mock.On("Convert", ctx, params)}
}

func (_c *CurrencyService_Convert_Call) Run(run func(ctx context.Context, params services.ConvertParams)) *CurrencyService_Convert_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(services.ConvertParams))
	})
	return _c
}

func (_c *CurrencyService_Convert_Call) Return(_a0 map[entities.Currency]float64, _a1 error) *CurrencyService_Convert_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CurrencyService_Convert_Call) RunAndReturn(run func(context.Context, services.ConvertParams) (map[entities.Currency]float64, error)) *CurrencyService_Convert_Call {
	_c.Call.Return(run)
	return _c
}

// NewCurrencyService creates a new instance of CurrencyService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCurrencyService(t interface {
	mock.TestingT
	Cleanup(func())
}) *CurrencyService {
	mock := &CurrencyService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}