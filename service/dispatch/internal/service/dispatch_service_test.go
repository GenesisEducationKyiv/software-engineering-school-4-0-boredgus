package service

// import (
// 	"context"
// 	"testing"
// 	"time"

// 	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/entities"
// 	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/mailing"
// 	client_mock "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/mocks/client"
// 	logger_mocks "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/mocks/logger"
// 	mailing_mocks "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/mocks/mailing"
// 	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/service/deps"

// 	repo_mocks "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/dispatch/internal/mocks/repo"
// 	"github.com/stretchr/testify/assert"
// )

// func Test_DispatchService_GetAllDispatches(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type mocked struct {
// 		dispatchesFromRepo  []deps.DispatchData
// 		getAllDispatchesErr error
// 	}

// 	dispatchRepoMock := repo_mocks.NewDispatchRepo(t)

// 	setup := func(m mocked, a args) func() {
// 		getCall := dispatchRepoMock.EXPECT().GetAllDispatches(a.ctx).
// 			Maybe().Return(m.dispatchesFromRepo, m.getAllDispatchesErr)

// 		return func() {
// 			getCall.Unset()
// 		}
// 	}

// 	dispatches := []deps.DispatchData{{
// 		Id:                 "id",
// 		Label:              "label",
// 		CountOfSubscribers: 2,
// 	}}
// 	ctx := context.Background()
// 	arguments := args{ctx: context.Background()}
// 	tests := []struct {
// 		name           string
// 		args           args
// 		mockedValues   mocked
// 		expectedResult []deps.DispatchData
// 		expectedErr    error
// 	}{
// 		{
// 			name:         "failed: got an error from GetAllDispatches",
// 			args:         arguments,
// 			mockedValues: mocked{getAllDispatchesErr: assert.AnError},
// 			expectedErr:  assert.AnError,
// 		},
// 		{
// 			name:           "successfuly got all dispatches",
// 			args:           arguments,
// 			mockedValues:   mocked{dispatchesFromRepo: dispatches},
// 			expectedResult: dispatches,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			cleanup := setup(tt.mockedValues, tt.args)
// 			defer cleanup()

// 			s := &dispatchService{
// 				dispatchRepo: dispatchRepoMock,
// 			}
// 			actualResult, actualErr := s.GetAllDispatches(ctx)

// 			assert.Equal(t, tt.expectedResult, actualResult)
// 			if tt.expectedErr != nil {
// 				assert.ErrorIs(t, actualErr, tt.expectedErr)

// 				return
// 			}
// 			assert.Nil(t, actualErr)
// 		})
// 	}
// }

// func Test_DispatchService_SubscribeForDispatch(t *testing.T) {
// 	type args struct {
// 		ctx        context.Context
// 		email      string
// 		dispatchId string
// 	}
// 	type mocked struct {
// 		getDispatchErr error
// 		createUserErr  error
// 		createSubErr   error
// 	}

// 	dispatchRepoMock := repo_mocks.NewDispatchRepo(t)
// 	userRepoMock := repo_mocks.NewUserRepo(t)
// 	subRepoMock := repo_mocks.NewSubRepo(t)
// 	setup := func(m *mocked, a *args) func() {
// 		getDsptchCall := dispatchRepoMock.EXPECT().
// 			GetDispatchByID(a.ctx, a.dispatchId).
// 			Maybe().Return(entities.CurrencyDispatch{}, m.getDispatchErr)
// 		createUserCall := userRepoMock.EXPECT().
// 			CreateUser(a.ctx, a.email).
// 			Maybe().NotBefore(getDsptchCall).Return(m.createUserErr)
// 		createSubCall := subRepoMock.EXPECT().
// 			CreateSubscription(a.ctx, deps.SubscriptionData{
// 				Email:    a.email,
// 				Dispatch: a.dispatchId,
// 			}).Maybe().NotBefore(createUserCall).Return(m.createSubErr)

// 		return func() {
// 			getDsptchCall.Unset()
// 			createUserCall.Unset()
// 			createSubCall.Unset()
// 		}
// 	}

// 	a := args{
// 		ctx:   context.Background(),
// 		email: "example@gmail.com",
// 	}

// 	tests := []struct {
// 		name         string
// 		mockedValues *mocked
// 		args         *args
// 		expectedErr  error
// 	}{
// 		{
// 			name:         "failed: got error from GetDispatchByID",
// 			args:         &a,
// 			mockedValues: &mocked{getDispatchErr: assert.AnError},
// 			expectedErr:  assert.AnError,
// 		},
// 		{
// 			name:         "failed: user already subscribed for this dispatch",
// 			args:         &a,
// 			mockedValues: &mocked{createUserErr: assert.AnError},
// 			expectedErr:  assert.AnError,
// 		},
// 		{
// 			name:         "failed: got an error from CreateSubscription",
// 			args:         &a,
// 			mockedValues: &mocked{createSubErr: assert.AnError},
// 			expectedErr:  assert.AnError,
// 		},
// 		{
// 			name:         "successfuly subscribed for a dispatch",
// 			args:         &a,
// 			mockedValues: &mocked{},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			cleanup := setup(tt.mockedValues, tt.args)
// 			defer cleanup()
// 			s := &dispatchService{
// 				userRepo:     userRepoMock,
// 				subRepo:      subRepoMock,
// 				dispatchRepo: dispatchRepoMock,
// 			}

// 			actualErr := s.SubscribeForDispatch(tt.args.ctx, tt.args.email, tt.args.dispatchId)
// 			if tt.expectedErr != nil {
// 				assert.ErrorIs(t, actualErr, tt.expectedErr)

// 				return
// 			}
// 			assert.Nil(t, actualErr)
// 		})
// 	}
// }

// func Test_DispatchService_SendDispatch(t *testing.T) {
// 	type args struct {
// 		ctx        context.Context
// 		dispatchId string
// 	}
// 	type mocked struct {
// 		dispatch       entities.CurrencyDispatch
// 		getDispatchErr error
// 		subscribers    []string
// 		getSubsErr     error
// 		rates          map[string]float64
// 		convertErr     error
// 		parsedEmail    []byte
// 		sendErr        error
// 	}

// 	loggerMock := logger_mocks.NewLogger()
// 	dispatchRepoMock := repo_mocks.NewDispatchRepo(t)
// 	csClientMock := client_mock.NewCurrencyServiceClient(t)
// 	mailmanMock := mailing_mocks.NewMailman(t)

// 	setup := func(m *mocked, a *args) func() {
// 		getDsptchCall := dispatchRepoMock.EXPECT().
// 			GetDispatchByID(a.ctx, a.dispatchId).
// 			Maybe().Return(m.dispatch, m.getDispatchErr)
// 		getSubsCall := dispatchRepoMock.EXPECT().
// 			GetSubscribersOfDispatch(a.ctx, a.dispatchId).
// 			Maybe().NotBefore(getDsptchCall).Return(m.subscribers, m.getSubsErr)
// 		convertCall := csClientMock.EXPECT().
// 			Convert(a.ctx, m.dispatch.Details.BaseCurrency, m.dispatch.Details.TargetCurrencies).Maybe().NotBefore(getSubsCall).Return(m.rates, m.convertErr)
// 		sendCall := mailmanMock.EXPECT().
// 			Send(mailing.Email{
// 				To:       m.subscribers,
// 				Subject:  m.dispatch.Label,
// 				HTMLBody: string(m.parsedEmail),
// 			}).Maybe().NotBefore(convertCall).Return(m.sendErr)

// 		return func() {
// 			getDsptchCall.Unset()
// 			getSubsCall.Unset()
// 			convertCall.Unset()
// 			sendCall.Unset()
// 		}
// 	}

// 	a := args{
// 		ctx:        context.Background(),
// 		dispatchId: "dispatch-id",
// 	}
// 	invalidDispatch := entities.CurrencyDispatch{
// 		Id:           "id",
// 		Label:        "label",
// 		SendAt:       time.Date(1, 1, 1, 1, 1, 1, 1, time.UTC),
// 		TemplateName: "template",
// 		Details: entities.CurrencyDispatchDetails{
// 			BaseCurrency:     "base",
// 			TargetCurrencies: []string{"target"}},
// 		CountOfSubscribers: 2,
// 	}
// 	dispatch := entities.CurrencyDispatch{
// 		Id:           "id",
// 		Label:        "label",
// 		SendAt:       time.Date(1, 1, 1, 1, 1, 1, 1, time.UTC),
// 		TemplateName: "test/test",
// 		Details: entities.CurrencyDispatchDetails{
// 			BaseCurrency:     "base",
// 			TargetCurrencies: []string{"target"}},
// 		CountOfSubscribers: 2,
// 	}
// 	subscribers := []string{"sub1", "sub2"}

// 	tests := []struct {
// 		name         string
// 		args         *args
// 		mockedValues *mocked
// 		wantErr      error
// 	}{
// 		{
// 			name:         "failed: got an error from GetDispatchByID",
// 			args:         &a,
// 			mockedValues: &mocked{getDispatchErr: assert.AnError},
// 			wantErr:      assert.AnError,
// 		},
// 		{
// 			name:         "success: there is no subscribers for dispatch",
// 			args:         &a,
// 			mockedValues: &mocked{},
// 		},
// 		{
// 			name: "failed: get an error from GetSubscribersOfDispatch",
// 			args: &a,
// 			mockedValues: &mocked{
// 				dispatch:   dispatch,
// 				getSubsErr: assert.AnError,
// 			},
// 			wantErr: assert.AnError,
// 		},
// 		{
// 			name: "failed: got an error from Convert",
// 			args: &a,
// 			mockedValues: &mocked{
// 				subscribers: subscribers,
// 				dispatch:    dispatch,
// 				convertErr:  assert.AnError,
// 			},
// 			wantErr: assert.AnError,
// 		},
// 		{
// 			name: "failed: got an error while parsing the template",
// 			args: &a,
// 			mockedValues: &mocked{
// 				subscribers: subscribers,
// 				dispatch:    invalidDispatch,
// 			},
// 			wantErr: TemplateParseErr,
// 		},
// 		{
// 			name: "failed: got an error from mailman",
// 			args: &a,
// 			mockedValues: &mocked{
// 				subscribers: subscribers,
// 				dispatch:    dispatch,
// 				parsedEmail: []byte("<div><span>test template</span></div>"),
// 				sendErr:     assert.AnError,
// 			},
// 			wantErr: assert.AnError,
// 		},
// 		{
// 			name: "successfuly sent",
// 			args: &a,
// 			mockedValues: &mocked{
// 				subscribers: subscribers,
// 				dispatch:    dispatch,
// 				parsedEmail: []byte("<div><span>test template</span></div>"),
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			cleanup := setup(tt.mockedValues, tt.args)
// 			defer cleanup()

// 			s := &dispatchService{
// 				dispatchRepo: dispatchRepoMock,
// 				mailman:      mailmanMock,
// 				csClient:     csClientMock,
// 				log:          loggerMock,
// 			}
// 			actualErr := s.SendDispatch(tt.args.ctx, tt.args.dispatchId)

// 			if tt.wantErr != nil {
// 				assert.ErrorIs(t, actualErr, tt.wantErr)

// 				return
// 			}
// 			assert.Nil(t, actualErr)
// 		})
// 	}
// }
