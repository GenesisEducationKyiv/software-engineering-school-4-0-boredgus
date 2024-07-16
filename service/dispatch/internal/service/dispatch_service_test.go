package service_test

import (
	"context"
	"testing"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/entities"
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

	a := args{
		ctx:        context.Background(),
		email:      "example@gmail.com",
		dispatchId: "dispatch-id",
	}

	dispatch := entities.CurrencyDispatch{
		ID: a.dispatchId,
	}

	tests := []struct {
		name                 string
		setup                func(*args) func()
		args                 *args
		expectedSubscription *entities.Subscription
		expectedErr          error
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
					Return(entities.CurrencyDispatch{}, nil)
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
					Return(entities.SubscriptionStatusActive, nil)

				return func() {
					getDispatchCall.Unset()
					getStatusCall.Unset()
				}
			},
			expectedErr: service.ErrAlreadyExists,
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
					Return(0, service.ErrNotFound)
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
					Return(0, service.ErrNotFound)
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
					Return(0, service.ErrNotFound)
				createUserCall := userRepoMock.EXPECT().
					CreateUser(a.ctx, a.email).Once().
					Return(nil)
				createSubscriptionCall := subRepoMock.EXPECT().
					CreateSubscription(a.ctx, service.SubscriptionData{Email: a.email, DispatchID: a.dispatchId}).
					Return(nil)

				return func() {
					getDispatchCall.Unset()
					getStatusCall.Unset()
					createUserCall.Unset()
					createSubscriptionCall.Unset()
				}
			},
			expectedSubscription: &entities.Subscription{
				DispatchID: a.dispatchId,
				Email:      a.email,
				Status:     entities.SubscriptionStatusActive,
			},
			expectedErr: nil,
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
					Return(entities.SubscriptionStatusCancelled, nil)
				updateStatusCall := subRepoMock.EXPECT().
					UpdateSubscriptionStatus(a.ctx,
						service.SubscriptionData{Email: a.email, DispatchID: a.dispatchId},
						entities.SubscriptionStatusActive).Return(assert.AnError)

				return func() {
					getDispatchCall.Unset()
					getStatusCall.Unset()
					updateStatusCall.Unset()
				}
			},
			expectedErr: assert.AnError,
		},
		{
			name: "success: reactivated subscription",
			args: &a,
			setup: func(a *args) func() {
				getDispatchCall := dispatchRepoMock.EXPECT().
					GetDispatchByID(a.ctx, a.dispatchId).Once().
					Return(dispatch, nil)
				getStatusCall := subRepoMock.EXPECT().
					GetStatusOfSubscription(a.ctx, service.SubscriptionData{Email: a.email, DispatchID: a.dispatchId}).
					Return(entities.SubscriptionStatusCancelled, nil)
				updateStatusCall := subRepoMock.EXPECT().
					UpdateSubscriptionStatus(a.ctx,
						service.SubscriptionData{Email: a.email, DispatchID: a.dispatchId},
						entities.SubscriptionStatusActive).Return(nil)

				return func() {
					getDispatchCall.Unset()
					getStatusCall.Unset()
					updateStatusCall.Unset()
				}
			},
			expectedSubscription: dispatch.ToSubscription(a.email, entities.SubscriptionStatusActive),
			expectedErr:          nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := tt.setup(tt.args)
			defer cleanup()

			s := service.NewDispatchService(userRepoMock, subRepoMock, dispatchRepoMock)
			actualSubscrioption, actualErr := s.SubscribeForDispatch(tt.args.ctx, tt.args.email, tt.args.dispatchId)

			assert.Equal(t, tt.expectedSubscription, actualSubscrioption)
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

	arguments := args{
		ctx:        context.Background(),
		email:      "email.@email.com",
		dispatchId: "dispatch-id",
	}

	tests := []struct {
		name                 string
		args                 args
		setup                func(args) func()
		expectedSubscription *entities.Subscription
		expectedErr          error
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
						entities.SubscriptionStatusCancelled).
					Return(service.ErrNotFound)

				return func() {
					getDispatchCall.Unset()
					updateStatusCall.Unset()
				}
			},
			expectedErr: service.ErrNotFound,
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
						entities.SubscriptionStatusCancelled).
					Return(nil)

				return func() {
					getDispatchCall.Unset()
					updateStatusCall.Unset()
				}
			},
			expectedSubscription: &entities.Subscription{
				DispatchID: arguments.dispatchId,
				Email:      arguments.email,
				Status:     entities.SubscriptionStatusCancelled,
			},
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := tt.setup(tt.args)
			defer cleanup()

			s := service.NewDispatchService(userRepoMock, subRepoMock, dispatchRepoMock)
			actualSubscription, actualErr := s.UnsubscribeFromDispatch(tt.args.ctx, tt.args.email, tt.args.dispatchId)

			assert.Equal(t, tt.expectedSubscription, actualSubscription)
			if tt.expectedErr != nil {
				assert.ErrorIs(t, actualErr, tt.expectedErr)

				return
			}
			assert.Nil(t, actualErr)
		})
	}
}
