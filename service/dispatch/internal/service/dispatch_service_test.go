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
			name: "failed: subscription is active",
			args: &a,
			setup: func(a *args) func() {
				getDispatchCall := dispatchRepoMock.EXPECT().
					GetDispatchByID(a.ctx, a.dispatchId).Once().
					Return(entities.CurrencyDispatch{ID: a.dispatchId}, nil)
				getStatusCall := subRepoMock.EXPECT().
					GetStatusOfSubscription(a.ctx, service.SubscriptionData{Email: a.email, DispatchID: a.dispatchId}).
					Return(service.SubscriptionRenewedStatus, nil)

				return func() {
					getDispatchCall.Unset()
					getStatusCall.Unset()
				}
			},
			expectedErr: service.AlreadyExistsErr,
		},
		{
			name: "failed: got an error from UpdateSubscriptionStatus",
			args: &a,
			setup: func(a *args) func() {
				getDispatchCall := dispatchRepoMock.EXPECT().
					GetDispatchByID(a.ctx, a.dispatchId).Once().
					Return(entities.CurrencyDispatch{ID: a.dispatchId}, nil)
				getStatusCall := subRepoMock.EXPECT().
					GetStatusOfSubscription(a.ctx, service.SubscriptionData{Email: a.email, DispatchID: a.dispatchId}).
					Return(service.SubscriptionCancelledStatus, nil)
				updateStatusCall := subRepoMock.EXPECT().
					UpdateSubscriptionStatus(a.ctx,
						service.SubscriptionData{Email: a.email, DispatchID: a.dispatchId},
						service.SubscriptionRenewedStatus).Return(assert.AnError)

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
				getDispatchCall := dispatchRepoMock.EXPECT().
					GetDispatchByID(a.ctx, a.dispatchId).Once().
					Return(dispatch, nil)
				getStatusCall := subRepoMock.EXPECT().
					GetStatusOfSubscription(a.ctx, service.SubscriptionData{Email: a.email, DispatchID: a.dispatchId}).
					Return(service.SubscriptionCancelledStatus, nil)
				updateStatusCall := subRepoMock.EXPECT().
					UpdateSubscriptionStatus(a.ctx,
						service.SubscriptionData{Email: a.email, DispatchID: a.dispatchId},
						service.SubscriptionRenewedStatus).Return(nil)
				brokerCall := brokerMock.EXPECT().
					RenewSubscription(service.DispatchToSubscription(dispatch, a.email))

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
				getStatusCall := subRepoMock.EXPECT().
					GetStatusOfSubscription(a.ctx, service.SubscriptionData{Email: a.email, DispatchID: a.dispatchId}).
					Return(0, service.NotFoundErr)

				return func() {
					getDispatchCall.Unset()
					getStatusCall.Unset()
				}
			},
			expectedErr: service.NotFoundErr,
		},
		{
			name: "failed: got unknown error from GetStatusOfSubscription",
			args: arguments,
			setup: func(a args) func() {
				dispatch := entities.CurrencyDispatch{ID: a.dispatchId}
				getDispatchCall := dispatchRepoMock.EXPECT().
					GetDispatchByID(a.ctx, a.dispatchId).Once().
					Return(dispatch, nil)
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
			name: "success: subcription is already cancelled",
			args: arguments,
			setup: func(a args) func() {
				dispatch := entities.CurrencyDispatch{ID: a.dispatchId}
				getDispatchCall := dispatchRepoMock.EXPECT().
					GetDispatchByID(a.ctx, a.dispatchId).Once().
					Return(dispatch, nil)
				getStatusCall := subRepoMock.EXPECT().
					GetStatusOfSubscription(a.ctx, service.SubscriptionData{Email: a.email, DispatchID: a.dispatchId}).
					Return(service.SubscriptionCancelledStatus, nil)

				return func() {
					getDispatchCall.Unset()
					getStatusCall.Unset()
				}
			},
			expectedErr: nil,
		},
		{
			name: "failed: got an error from UpdateSubscriptionStatus",
			args: arguments,
			setup: func(a args) func() {
				dispatch := entities.CurrencyDispatch{ID: a.dispatchId}
				getDispatchCall := dispatchRepoMock.EXPECT().
					GetDispatchByID(a.ctx, a.dispatchId).Once().
					Return(dispatch, nil)
				getStatusCall := subRepoMock.EXPECT().
					GetStatusOfSubscription(a.ctx, service.SubscriptionData{Email: a.email, DispatchID: a.dispatchId}).
					Return(service.SubscriptionRenewedStatus, nil).Once()
				updateStatusCall := subRepoMock.EXPECT().
					UpdateSubscriptionStatus(a.ctx,
						service.SubscriptionData{Email: a.email, DispatchID: a.dispatchId},
						service.SubscriptionCancelledStatus).Return(assert.AnError).Once()

				return func() {
					getDispatchCall.Unset()
					getStatusCall.Unset()
					updateStatusCall.Unset()
				}
			},
			expectedErr: assert.AnError,
		},
		{
			name: "success: cancelld subscription",
			args: arguments,
			setup: func(a args) func() {
				dispatch := entities.CurrencyDispatch{ID: a.dispatchId}
				getDispatchCall := dispatchRepoMock.EXPECT().
					GetDispatchByID(a.ctx, a.dispatchId).Once().
					Return(dispatch, nil)
				getStatusCall := subRepoMock.EXPECT().
					GetStatusOfSubscription(a.ctx, service.SubscriptionData{Email: a.email, DispatchID: a.dispatchId}).
					Return(service.SubscriptionCreatedStatus, nil).Once()
				updateStatusCall := subRepoMock.EXPECT().
					UpdateSubscriptionStatus(a.ctx,
						service.SubscriptionData{Email: a.email, DispatchID: a.dispatchId},
						service.SubscriptionCancelledStatus).Return(nil).Once()
				brokerCall := brokerMock.EXPECT().
					CancelSubscription(service.DispatchToSubscription(dispatch, a.email))

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
			actualErr := s.UnsubscribeFromDispatch(tt.args.ctx, tt.args.email, tt.args.dispatchId)

			if tt.expectedErr != nil {
				assert.ErrorIs(t, actualErr, tt.expectedErr)

				return
			}
			assert.Nil(t, actualErr)
		})
	}
}
