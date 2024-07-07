package service_test

import (
	"context"
	"testing"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/entities"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/service"

	broker_mock "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/mocks/broker"
	repo_mocks "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/mocks/repo"
	"github.com/stretchr/testify/assert"
)

func Test_DispatchService_GetAllDispatches(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type mocked struct {
		dispatchesFromRepo  []service.DispatchData
		getAllDispatchesErr error
	}

	dispatchRepoMock := repo_mocks.NewDispatchRepo(t)

	setup := func(m mocked, a args) func() {
		getCall := dispatchRepoMock.EXPECT().GetAllDispatches(a.ctx).
			Maybe().Return(m.dispatchesFromRepo, m.getAllDispatchesErr)

		return func() {
			getCall.Unset()
		}
	}

	dispatches := []service.DispatchData{{
		Id:                 "id",
		Label:              "label",
		CountOfSubscribers: 2,
	}}
	ctx := context.Background()
	arguments := args{ctx: context.Background()}
	tests := []struct {
		name           string
		args           args
		mockedValues   mocked
		expectedResult []service.DispatchData
		expectedErr    error
	}{
		{
			name:         "failed: got an error from GetAllDispatches",
			args:         arguments,
			mockedValues: mocked{getAllDispatchesErr: assert.AnError},
			expectedErr:  assert.AnError,
		},
		{
			name:           "successfuly got all dispatches",
			args:           arguments,
			mockedValues:   mocked{dispatchesFromRepo: dispatches},
			expectedResult: dispatches,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setup(tt.mockedValues, tt.args)
			defer cleanup()

			s := service.NewDispatchService(nil, nil, nil, nil, nil, dispatchRepoMock)
			actualResult, actualErr := s.GetAllDispatches(ctx)

			assert.Equal(t, tt.expectedResult, actualResult)
			if tt.expectedErr != nil {
				assert.ErrorIs(t, actualErr, tt.expectedErr)

				return
			}
			assert.Nil(t, actualErr)
		})
	}
}

func Test_DispatchService_SubscribeForDispatch(t *testing.T) {
	type args struct {
		ctx        context.Context
		email      string
		dispatchId string
	}
	type mocked struct {
		dispatch       entities.CurrencyDispatch
		getDispatchErr error
		createUserErr  error
		createSubErr   error
	}

	dispatchRepoMock := repo_mocks.NewDispatchRepo(t)
	userRepoMock := repo_mocks.NewUserRepo(t)
	subRepoMock := repo_mocks.NewSubRepo(t)
	brokerMock := broker_mock.NewBroker(t)

	setup := func(m *mocked, a *args) func() {
		getDsptchCall := dispatchRepoMock.EXPECT().
			GetDispatchByID(a.ctx, a.dispatchId).
			Maybe().Return(m.dispatch, m.getDispatchErr)
		createUserCall := userRepoMock.EXPECT().
			CreateUser(a.ctx, a.email).
			Maybe().NotBefore(getDsptchCall).Return(m.createUserErr)
		createSubCall := subRepoMock.EXPECT().
			CreateSubscription(a.ctx, service.SubscriptionData{
				Email:    a.email,
				Dispatch: a.dispatchId,
			}).Maybe().NotBefore(createUserCall).Return(m.createSubErr)
		brokerCall := brokerMock.EXPECT().
			CreateSubscription(service.Subscription{
				DispatchID:  m.dispatch.Id,
				Email:       a.email,
				BaseCcy:     m.dispatch.Details.BaseCurrency,
				TargetCcies: m.dispatch.Details.TargetCurrencies,
				SendAt:      m.dispatch.SendAt,
			}).Maybe().NotBefore(createSubCall)

		return func() {
			getDsptchCall.Unset()
			createUserCall.Unset()
			createSubCall.Unset()
			brokerCall.Unset()
		}
	}

	a := args{
		ctx:   context.Background(),
		email: "example@gmail.com",
	}

	tests := []struct {
		name         string
		mockedValues *mocked
		args         *args
		expectedErr  error
	}{
		{
			name:         "failed: got error from GetDispatchByID",
			args:         &a,
			mockedValues: &mocked{getDispatchErr: assert.AnError},
			expectedErr:  assert.AnError,
		},
		{
			name:         "failed: failed to create user",
			args:         &a,
			mockedValues: &mocked{createUserErr: assert.AnError},
			expectedErr:  assert.AnError,
		},
		{
			name:         "failed: got an error from CreateSubscription",
			args:         &a,
			mockedValues: &mocked{createSubErr: assert.AnError},
			expectedErr:  assert.AnError,
		},
		{
			name:         "successfuly subscribed for a dispatch",
			args:         &a,
			mockedValues: &mocked{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setup(tt.mockedValues, tt.args)
			defer cleanup()

			s := service.NewDispatchService(nil, nil, nil, userRepoMock, subRepoMock, dispatchRepoMock)
			actualErr := s.SubscribeForDispatch(tt.args.ctx, tt.args.email, tt.args.dispatchId)

			if tt.expectedErr != nil {
				assert.ErrorIs(t, actualErr, tt.expectedErr)

				return
			}
			assert.Nil(t, actualErr)
		})
	}
}
