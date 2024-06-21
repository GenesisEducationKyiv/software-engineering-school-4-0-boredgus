package dispatch_service

import (
	"context"
	"subscription-api/internal/db"
	"subscription-api/internal/entities"
	e "subscription-api/internal/entities"
	"subscription-api/internal/mailing"
	client_mocks "subscription-api/internal/mocks/clients"
	config_mocks "subscription-api/internal/mocks/config"
	db_mocks "subscription-api/internal/mocks/db"
	mailing_mocks "subscription-api/internal/mocks/mailing"
	repo_mocks "subscription-api/internal/mocks/repo"
	"subscription-api/internal/services"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_DispatchService_GetAllDispatches(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type mocked struct {
		dispatchesFromRepo  []services.DispatchData
		getAllDispatchesErr error
	}

	storeMock := db_mocks.NewStore()
	dispatchRepoMock := repo_mocks.NewDispatchRepo(t)

	setup := func(m mocked, a args) func() {
		getCall := dispatchRepoMock.EXPECT().GetAllDispatches(a.ctx, mock.Anything).
			Maybe().Return(m.dispatchesFromRepo, m.getAllDispatchesErr)

		return func() {
			getCall.Unset()
		}
	}

	dispatches := []services.DispatchData{{
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
		expectedResult []services.DispatchData
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

			s := &dispatchService{
				store:        storeMock,
				dispatchRepo: dispatchRepoMock,
			}
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
		getDispatchErr error
		createUserErr  error
		createSubErr   error
	}

	storeMock := db_mocks.NewStore()
	dispatchRepoMock := repo_mocks.NewDispatchRepo(t)
	userRepoMock := repo_mocks.NewUserRepo(t)
	subRepoMock := repo_mocks.NewSubRepo(t)
	setup := func(m *mocked, a *args) func() {
		getDsptchCall := dispatchRepoMock.EXPECT().
			GetDispatchByID(a.ctx, mock.Anything, a.dispatchId).
			Maybe().Return(e.CurrencyDispatch{}, m.getDispatchErr)
		createUserCall := userRepoMock.EXPECT().
			CreateUser(a.ctx, mock.Anything, a.email).
			Maybe().NotBefore(getDsptchCall).Return(m.createUserErr)
		createSubCall := subRepoMock.EXPECT().
			CreateSubscription(a.ctx, mock.Anything, db.SubscriptionData{
				Email: a.email,
			}).Maybe().NotBefore(createUserCall).Return(m.createSubErr)

		return func() {
			getDsptchCall.Unset()
			createUserCall.Unset()
			createSubCall.Unset()
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
			name:         "failed: user already subscribed for this dispatch",
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
			s := &dispatchService{
				store:        storeMock,
				userRepo:     userRepoMock,
				subRepo:      subRepoMock,
				dispatchRepo: dispatchRepoMock,
			}

			actualErr := s.SubscribeForDispatch(tt.args.ctx, tt.args.email, tt.args.dispatchId)
			if tt.expectedErr != nil {
				assert.ErrorIs(t, actualErr, tt.expectedErr)

				return
			}
			assert.Nil(t, actualErr)
		})
	}
}

func Test_DispatchService_SendDispatch(t *testing.T) {
	type args struct {
		ctx        context.Context
		dispatchId string
	}
	type mocked struct {
		dispatch       entities.CurrencyDispatch
		getDispatchErr error
		subscribers    []string
		getSubsErr     error
		rates          map[string]float64
		convertErr     error
		parsedEmail    []byte
		sendErr        error
	}

	storeMock := db_mocks.NewStore()
	loggerMock := config_mocks.NewLogger()
	dispatchRepoMock := repo_mocks.NewDispatchRepo(t)
	csClientMock := client_mocks.NewCurrencyServiceClient(t)
	mailmanMock := mailing_mocks.NewMailman(t)

	setup := func(m *mocked, a *args) func() {
		getDsptchCall := dispatchRepoMock.EXPECT().
			GetDispatchByID(a.ctx, mock.Anything, a.dispatchId).
			Maybe().Return(m.dispatch, m.getDispatchErr)
		getSubsCall := dispatchRepoMock.EXPECT().
			GetSubscribersOfDispatch(a.ctx, mock.Anything, a.dispatchId).
			Maybe().NotBefore(getDsptchCall).Return(m.subscribers, m.getSubsErr)
		convertCall := csClientMock.EXPECT().
			Convert(a.ctx, services.ConvertCurrencyParams{
				Base:   m.dispatch.Details.BaseCurrency,
				Target: m.dispatch.Details.TargetCurrencies,
			}).Maybe().NotBefore(getSubsCall).Return(m.rates, m.convertErr)
		sendCall := mailmanMock.EXPECT().
			Send(mailing.Email{
				To:       m.subscribers,
				Subject:  m.dispatch.Label,
				HTMLBody: string(m.parsedEmail),
			}).Maybe().NotBefore(convertCall).Return(m.sendErr)

		return func() {
			getDsptchCall.Unset()
			getSubsCall.Unset()
			convertCall.Unset()
			sendCall.Unset()
		}
	}

	a := args{
		ctx:        context.Background(),
		dispatchId: "dispatch-id",
	}
	invalidDispatch := e.CurrencyDispatch{
		Id:           "id",
		Label:        "label",
		SendAt:       time.Date(1, 1, 1, 1, 1, 1, 1, time.UTC),
		TemplateName: "template",
		Details: e.CurrencyDispatchDetails{
			BaseCurrency:     "base",
			TargetCurrencies: []string{"target"}},
		CountOfSubscribers: 2,
	}
	dispatch := e.CurrencyDispatch{
		Id:           "id",
		Label:        "label",
		SendAt:       time.Date(1, 1, 1, 1, 1, 1, 1, time.UTC),
		TemplateName: "test/test",
		Details: e.CurrencyDispatchDetails{
			BaseCurrency:     "base",
			TargetCurrencies: []string{"target"}},
		CountOfSubscribers: 2,
	}
	subscribers := []string{"sub1", "sub2"}

	tests := []struct {
		name         string
		args         *args
		mockedValues *mocked
		wantErr      error
	}{
		{
			name:         "failed: got an error from GetDispatchByID",
			args:         &a,
			mockedValues: &mocked{getDispatchErr: assert.AnError},
			wantErr:      assert.AnError,
		},
		{
			name:         "failed: get an error from GetSubscribersOfDispatch",
			args:         &a,
			mockedValues: &mocked{getSubsErr: assert.AnError},
			wantErr:      assert.AnError,
		},
		{
			name:         "success: there is no subscribers for dispatch",
			args:         &a,
			mockedValues: &mocked{},
		},
		{
			name: "failed: got an error from Convert",
			args: &a,
			mockedValues: &mocked{
				subscribers: subscribers,
				dispatch:    invalidDispatch,
				convertErr:  assert.AnError,
			},
			wantErr: assert.AnError,
		},
		{
			name: "failed: got an error while parsing the template",
			args: &a,
			mockedValues: &mocked{
				subscribers: subscribers,
				dispatch:    invalidDispatch,
			},
			wantErr: TemplateParseErr,
		},
		{
			name: "failed: got an error from mailman",
			args: &a,
			mockedValues: &mocked{
				subscribers: subscribers,
				dispatch:    dispatch,
				parsedEmail: []byte("<div><span>test template</span></div>"),
				sendErr:     assert.AnError,
			},
			wantErr: assert.AnError,
		},
		{
			name: "successfuly sent",
			args: &a,
			mockedValues: &mocked{
				subscribers: subscribers,
				dispatch:    dispatch,
				parsedEmail: []byte("<div><span>test template</span></div>"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setup(tt.mockedValues, tt.args)
			defer cleanup()

			s := &dispatchService{
				store:        storeMock,
				dispatchRepo: dispatchRepoMock,
				mailman:      mailmanMock,
				csClient:     csClientMock,
				log:          loggerMock,
			}
			actualErr := s.SendDispatch(tt.args.ctx, tt.args.dispatchId)

			if tt.wantErr != nil {
				assert.ErrorIs(t, actualErr, tt.wantErr)

				return
			}
			assert.Nil(t, actualErr)
		})
	}
}
