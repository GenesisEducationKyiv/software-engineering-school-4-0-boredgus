package service_test

import (
	"context"
	"testing"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/entities"
	broker_mock "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/mocks/broker"
	repo_mocks "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/mocks/repo"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/service"
	"github.com/stretchr/testify/assert"
)

func Test_DispatchService_SubscribeForDispatch(t *testing.T) {
	type args struct {
		ctx        context.Context
		email      string
		dispatchId string
	}

	dispatchRepoMock := repo_mocks.NewDispatchRepo(t)
	userRepoMock := repo_mocks.NewUserRepo(t)
	subRepoMock := repo_mocks.NewSubRepo(t)
	brokerMock := broker_mock.NewBroker(t)

	a := args{
		ctx:        context.Background(),
		email:      "example@gmail.com",
		dispatchId: "dispatch-id",
	}

	tests := []struct {
		name        string
		setup       func(*args) func()
		args        *args
		expectedErr error
	}{
		{
			name: "failed: got error from GetDispatchByID",
			args: &a,
			setup: func(a *args) func() {
				getDispatchCall := dispatchRepoMock.EXPECT().
					GetDispatchByID(a.ctx, a.dispatchId).Once().
					Return(entities.CurrencyDispatch{}, assert.AnError)

				return func() {
					getDispatchCall.Unset()
				}
			},
			expectedErr: assert.AnError,
		},
		{
			name: "failed: got an error from GetStatusOfSubscription",
			args: &a,
			setup: func(a *args) func() {
				getDispatchCall := dispatchRepoMock.EXPECT().
					GetDispatchByID(a.ctx, a.dispatchId).Once().
					Return(entities.CurrencyDispatch{}, assert.AnError)
				getStatusCall := subRepoMock.EXPECT().
					GetStatusOfSubscription(a.ctx, service.SubscriptionData{Email: a.email, DispatchID: a.dispatchId}).
					Return(0, assert.AnError)

				return func() {
					getDispatchCall.Unset()
					getStatusCall.Unset()
				}
			},
			expectedErr: assert.AnError,
		},
		{
			name: "failed: subscription already exists",
			args: &a,
			setup: func(a *args) func() {
				getDispatchCall := dispatchRepoMock.EXPECT().
					GetDispatchByID(a.ctx, a.dispatchId).Once().
					Return(entities.CurrencyDispatch{}, nil)
				getStatusCall := subRepoMock.EXPECT().
					GetStatusOfSubscription(a.ctx, service.SubscriptionData{Email: a.email, DispatchID: a.dispatchId}).
					Return(service.SubscriptionStatusActive, nil)

				return func() {
					getDispatchCall.Unset()
					getStatusCall.Unset()
				}
			},
			expectedErr: service.AlreadyExistsErr,
		},
		{
			name: "failed: got an error on create user attempt",
			args: &a,
			setup: func(a *args) func() {
				getDispatchCall := dispatchRepoMock.EXPECT().
					GetDispatchByID(a.ctx, a.dispatchId).Once().
					Return(entities.CurrencyDispatch{ID: a.dispatchId}, nil)
				getStatusCall := subRepoMock.EXPECT().
					GetStatusOfSubscription(a.ctx, service.SubscriptionData{Email: a.email, DispatchID: a.dispatchId}).
					Return(0, service.NotFoundErr)
				createUserCall := userRepoMock.EXPECT().
					CreateUser(a.ctx, a.email).Once().
					Return(assert.AnError)

				return func() {
					getDispatchCall.Unset()
					getStatusCall.Unset()
					createUserCall.Unset()
				}
			},
			expectedErr: assert.AnError,
		},
		{
			name: "failed: got an error on create user attempt",
			args: &a,
			setup: func(a *args) func() {
				getDispatchCall := dispatchRepoMock.EXPECT().
					GetDispatchByID(a.ctx, a.dispatchId).Once().
					Return(entities.CurrencyDispatch{ID: a.dispatchId}, nil)
				getStatusCall := subRepoMock.EXPECT().
					GetStatusOfSubscription(a.ctx, service.SubscriptionData{Email: a.email, DispatchID: a.dispatchId}).
					Return(0, service.NotFoundErr)
				createUserCall := userRepoMock.EXPECT().
					CreateUser(a.ctx, a.email).Once().
					Return(assert.AnError)

				return func() {
					getDispatchCall.Unset()
					getStatusCall.Unset()
					createUserCall.Unset()
				}
			},
			expectedErr: assert.AnError,
		},
		{
			name: "failed: got an error on create subscription attempt",
			args: &a,
			setup: func(a *args) func() {
				getDispatchCall := dispatchRepoMock.EXPECT().
					GetDispatchByID(a.ctx, a.dispatchId).Once().
					Return(entities.CurrencyDispatch{ID: a.dispatchId}, nil)
				getStatusCall := subRepoMock.EXPECT().
					GetStatusOfSubscription(a.ctx, service.SubscriptionData{Email: a.email, DispatchID: a.dispatchId}).
					Return(0, service.NotFoundErr)
				createUserCall := userRepoMock.EXPECT().
					CreateUser(a.ctx, a.email).Once().
					Return(nil)
				createSubscriptionCall := subRepoMock.EXPECT().
					CreateSubscription(a.ctx, service.SubscriptionData{Email: a.email, DispatchID: a.dispatchId}).
					Return(assert.AnError).Once()

				return func() {
					getDispatchCall.Unset()
					getStatusCall.Unset()
					createUserCall.Unset()
					createSubscriptionCall.Unset()
				}
			},
			expectedErr: assert.AnError,
		},
		{
			name: "success: created new subscription",
			args: &a,
			setup: func(a *args) func() {
				dispatch := entities.CurrencyDispatch{ID: a.dispatchId}
				getDispatchCall := dispatchRepoMock.EXPECT().
					GetDispatchByID(a.ctx, a.dispatchId).Once().
					Return(dispatch, nil)
				getStatusCall := subRepoMock.EXPECT().
					GetStatusOfSubscription(a.ctx, service.SubscriptionData{Email: a.email, DispatchID: a.dispatchId}).
					Return(0, service.NotFoundErr)
				createUserCall := userRepoMock.EXPECT().
					CreateUser(a.ctx, a.email).Once().
					Return(nil)
				createSubscriptionCall := subRepoMock.EXPECT().
					CreateSubscription(a.ctx, service.SubscriptionData{Email: a.email, DispatchID: a.dispatchId}).
					Return(nil)
				brokerCall := brokerMock.EXPECT().
					CreateSubscription(service.DispatchToSubscription(dispatch, a.email))

				return func() {
					getDispatchCall.Unset()
					getStatusCall.Unset()
					createUserCall.Unset()
					createSubscriptionCall.Unset()
					brokerCall.Unset()
				}
			},
			expectedErr: nil,
		},
		{
			name: "failed: subscription already exists",
			args: &a,
			setup: func(a *args) func() {
				getDispatchCall := dispatchRepoMock.EXPECT().
					GetDispatchByID(a.ctx, a.dispatchId).Once().
					Return(entities.CurrencyDispatch{ID: a.dispatchId}, nil)
				getStatusCall := subRepoMock.EXPECT().
					GetStatusOfSubscription(a.ctx, service.SubscriptionData{Email: a.email, DispatchID: a.dispatchId}).
					Return(service.SubscriptionStatusActive, nil)

				return func() {
					getDispatchCall.Unset()
					getStatusCall.Unset()
				}
			},
			expectedErr: service.AlreadyExistsErr,
		},
		{
			name: "failed: got unknown error from GetStatusOfSubscription",
			args: &a,
			setup: func(a *args) func() {
				getDispatchCall := dispatchRepoMock.EXPECT().
					GetDispatchByID(a.ctx, a.dispatchId).Once().
					Return(entities.CurrencyDispatch{ID: a.dispatchId}, nil)
				getStatusCall := subRepoMock.EXPECT().
					GetStatusOfSubscription(a.ctx, service.SubscriptionData{Email: a.email, DispatchID: a.dispatchId}).
					Return(0, assert.AnError)

				return func() {
					getDispatchCall.Unset()
					getStatusCall.Unset()
				}
			},
			expectedErr: assert.AnError,
		},
		{
			name: "failed: got unknown error from GetStatusOfSubscription",
			args: &a,
			setup: func(a *args) func() {
				subData := service.SubscriptionData{Email: a.email, DispatchID: a.dispatchId}
				getDispatchCall := dispatchRepoMock.EXPECT().
					GetDispatchByID(a.ctx, a.dispatchId).Once().
					Return(entities.CurrencyDispatch{ID: a.dispatchId}, nil)
				getStatusCall := subRepoMock.EXPECT().
					GetStatusOfSubscription(a.ctx, subData).
					Return(service.SubscriptionStatusCancelled, nil)
				updateStatusCall := subRepoMock.EXPECT().
					UpdateSubscriptionStatus(a.ctx, subData, service.SubscriptionStatusActive).
					Return(assert.AnError)

				return func() {
					getDispatchCall.Unset()
					getStatusCall.Unset()
					updateStatusCall.Unset()
				}
			},
			expectedErr: assert.AnError,
		},
		{
			name: "success: renewed subscription",
			args: &a,
			setup: func(a *args) func() {
				dispatch := entities.CurrencyDispatch{ID: a.dispatchId}
				subData := service.SubscriptionData{Email: a.email, DispatchID: a.dispatchId}
				getDispatchCall := dispatchRepoMock.EXPECT().
					GetDispatchByID(a.ctx, a.dispatchId).Once().
					Return(entities.CurrencyDispatch{ID: a.dispatchId}, nil)
				getStatusCall := subRepoMock.EXPECT().
					GetStatusOfSubscription(a.ctx, subData).
					Return(service.SubscriptionStatusCancelled, nil)
				updateStatusCall := subRepoMock.EXPECT().
					UpdateSubscriptionStatus(a.ctx, subData, service.SubscriptionStatusActive).
					Return(nil)
				brokerCall := brokerMock.EXPECT().
					CreateSubscription(service.DispatchToSubscription(dispatch, a.email))

				return func() {
					getDispatchCall.Unset()
					getStatusCall.Unset()
					updateStatusCall.Unset()
					brokerCall.Unset()
				}
			},
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := tt.setup(tt.args)
			defer cleanup()

			s := service.NewDispatchService(userRepoMock, subRepoMock, dispatchRepoMock, brokerMock)
			actualErr := s.SubscribeForDispatch(tt.args.ctx, tt.args.email, tt.args.dispatchId)

			if tt.expectedErr != nil {
				assert.ErrorIs(t, actualErr, tt.expectedErr)

				return
			}
			assert.Nil(t, actualErr)
		})
	}
}

func Test_dispatchService_UnsubscribeFromDispatch(t *testing.T) {
	type args struct {
		ctx        context.Context
		email      string
		dispatchId string
	}

	dispatchRepoMock := repo_mocks.NewDispatchRepo(t)
	userRepoMock := repo_mocks.NewUserRepo(t)
	subRepoMock := repo_mocks.NewSubRepo(t)
	brokerMock := broker_mock.NewBroker(t)

	arguments := args{
		ctx:        context.Background(),
		email:      "email.@email.com",
		dispatchId: "dispatch-id",
	}

	tests := []struct {
		name        string
		args        args
		setup       func(args) func()
		expectedErr error
	}{
		{
			name: "failed: got an error from GetDispatchByID",
			args: arguments,
			setup: func(a args) func() {
				getDispatchCall := dispatchRepoMock.EXPECT().
					GetDispatchByID(a.ctx, a.dispatchId).Once().
					Return(entities.CurrencyDispatch{}, assert.AnError)

				return func() {
					getDispatchCall.Unset()
				}
			},
			expectedErr: assert.AnError,
		},
		{
			name: "failed: user is not subscribed for provided dispatch",
			args: arguments,
			setup: func(a args) func() {
				dispatch := entities.CurrencyDispatch{ID: a.dispatchId}
				getDispatchCall := dispatchRepoMock.EXPECT().
					GetDispatchByID(a.ctx, a.dispatchId).Once().
					Return(dispatch, nil)
				updateStatusCall := subRepoMock.EXPECT().
					UpdateSubscriptionStatus(a.ctx,
						service.SubscriptionData{Email: a.email, DispatchID: a.dispatchId},
						service.SubscriptionStatusCancelled).
					Return(service.NotFoundErr)

				return func() {
					getDispatchCall.Unset()
					updateStatusCall.Unset()
				}
			},
			expectedErr: service.NotFoundErr,
		},
		{
			name: "success: subcription cancelled",
			args: arguments,
			setup: func(a args) func() {
				dispatch := entities.CurrencyDispatch{ID: a.dispatchId}
				getDispatchCall := dispatchRepoMock.EXPECT().
					GetDispatchByID(a.ctx, a.dispatchId).Once().
					Return(dispatch, nil)
				updateStatusCall := subRepoMock.EXPECT().
					UpdateSubscriptionStatus(a.ctx,
						service.SubscriptionData{Email: a.email, DispatchID: a.dispatchId},
						service.SubscriptionStatusCancelled).
					Return(nil)
				brokerCall := brokerMock.EXPECT().
					CancelSubscription(service.DispatchToSubscription(dispatch, a.email))

				return func() {
					getDispatchCall.Unset()
					updateStatusCall.Unset()
					brokerCall.Unset()
				}
			},
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := tt.setup(tt.args)
			defer cleanup()

			s := service.NewDispatchService(userRepoMock, subRepoMock, dispatchRepoMock, brokerMock)
			actualErr := s.UnsubscribeFromDispatch(tt.args.ctx, tt.args.email, tt.args.dispatchId)

			if tt.expectedErr != nil {
				assert.ErrorIs(t, actualErr, tt.expectedErr)

				return
			}
			assert.Nil(t, actualErr)
		})
	}
}
