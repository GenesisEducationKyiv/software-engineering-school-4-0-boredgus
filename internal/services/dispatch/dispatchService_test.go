package dispatch_service

import (
	"context"
	"fmt"
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
		dispatches    []services.DispatchData
		getDsptchsErr error
	}

	storeMock := db_mocks.NewStore()
	dispatchRepoMock := repo_mocks.NewDispatchRepo(t)

	setup := func(m mocked, a args) func() {
		getCall := dispatchRepoMock.EXPECT().GetAllDispatches(a.ctx, mock.Anything).
			Maybe().Return(m.dispatches, m.getDsptchsErr)

		return func() {
			getCall.Unset()
		}
	}

	dispatches := []services.DispatchData{{
		Id:                 "id",
		Label:              "label",
		CountOfSubscribers: 2,
	}}
	someErr := fmt.Errorf("some err")
	ctx := context.Background()
	arguments := args{ctx: context.Background()}
	tests := []struct {
		name    string
		args    args
		mocked  mocked
		want    []services.DispatchData
		wantErr error
	}{
		{
			name:    "failed to get dispatches",
			args:    arguments,
			mocked:  mocked{getDsptchsErr: someErr},
			wantErr: someErr,
		},
		{
			name:   "success",
			args:   arguments,
			mocked: mocked{dispatches: dispatches},
			want:   dispatches,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setup(tt.mocked, tt.args)
			defer cleanup()

			s := &dispatchService{
				store:        storeMock,
				dispatchRepo: dispatchRepoMock,
			}
			got, err := s.GetAllDispatches(ctx)

			assert.Equal(t, tt.want, got)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)

				return
			}
			assert.Nil(t, err)
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
	getDispatchErr := fmt.Errorf("get-dispatch-err")
	createUserErr := fmt.Errorf("create-user-err")
	createSubErr := fmt.Errorf("create-sub-err")

	tests := []struct {
		name    string
		mocked  *mocked
		args    *args
		wantErr error
	}{
		{
			name:    "failed to get dispatch",
			args:    &a,
			mocked:  &mocked{getDispatchErr: getDispatchErr},
			wantErr: getDispatchErr,
		},
		{
			name:    "user already subscribed for such dispatch",
			args:    &a,
			mocked:  &mocked{createUserErr: createUserErr},
			wantErr: createUserErr,
		},
		{
			name:    "failed to create subscription",
			args:    &a,
			mocked:  &mocked{createSubErr: createSubErr},
			wantErr: createSubErr,
		},
		{
			name:   "success",
			args:   &a,
			mocked: &mocked{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setup(tt.mocked, tt.args)
			defer cleanup()
			s := &dispatchService{
				store:        storeMock,
				userRepo:     userRepoMock,
				subRepo:      subRepoMock,
				dispatchRepo: dispatchRepoMock,
			}

			err := s.SubscribeForDispatch(tt.args.ctx, tt.args.email, tt.args.dispatchId)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)

				return
			}
			assert.Nil(t, err)
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
	loggerMock := config_mocks.NewLogger(t)
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
		errorCall := loggerMock.EXPECT().
			Errorf(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()

		return func() {
			getDsptchCall.Unset()
			getSubsCall.Unset()
			convertCall.Unset()
			sendCall.Unset()
			errorCall.Unset()
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
	getDispatchErr := fmt.Errorf("get-dispatch-err")
	getSubsErr := fmt.Errorf("get-subs-err")
	subscribers := []string{"sub1", "sub2"}
	convertErr := fmt.Errorf("convert-err")
	sendErr := fmt.Errorf("send-err")

	tests := []struct {
		name    string
		args    *args
		mocked  *mocked
		wantErr error
	}{
		{
			name:    "failed to get dispatch data",
			args:    &a,
			mocked:  &mocked{getDispatchErr: getDispatchErr},
			wantErr: getDispatchErr,
		},
		{
			name: "failed to get subscribers of dispatch",
			args: &a,
			mocked: &mocked{
				getSubsErr: getSubsErr,
			},
			wantErr: getSubsErr,
		},
		{
			name:   "there is no subscribers for dispatch",
			args:   &a,
			mocked: &mocked{},
		},
		{
			name: "failed to convert currencies",
			args: &a,
			mocked: &mocked{
				subscribers: subscribers,
				dispatch:    invalidDispatch,
				convertErr:  convertErr,
			},
			wantErr: convertErr,
		},
		{
			name: "failed to parse email template",
			args: &a,
			mocked: &mocked{
				subscribers: subscribers,
				dispatch:    invalidDispatch,
			},
			wantErr: TemplateParseErr,
		},
		{
			name: "failed to send emails",
			args: &a,
			mocked: &mocked{
				subscribers: subscribers,
				dispatch:    dispatch,
				parsedEmail: []byte("<div><span>test template</span></div>"),
				sendErr:     sendErr,
			},
			wantErr: sendErr,
		},
		{
			name: "success",
			args: &a,
			mocked: &mocked{
				subscribers: subscribers,
				dispatch:    dispatch,
				parsedEmail: []byte("<div><span>test template</span></div>"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setup(tt.mocked, tt.args)
			defer cleanup()

			s := &dispatchService{
				store:        storeMock,
				dispatchRepo: dispatchRepoMock,
				mailman:      mailmanMock,
				csClient:     csClientMock,
				log:          loggerMock,
			}
			err := s.SendDispatch(tt.args.ctx, tt.args.dispatchId)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)

				return
			}
			assert.Nil(t, err)
		})
	}
}
