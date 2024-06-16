package ds

import (
	"context"
	"fmt"
	"subscription-api/internal/db"
	e "subscription-api/internal/entities"
	db_mocks "subscription-api/internal/mocks/db"
	repo_mocks "subscription-api/internal/mocks/repo"
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
		dispatches []e.CurrencyDispatch
		getByidErr error
	}

	storeMock := db_mocks.NewStore()
	dispatchRepoMock := repo_mocks.NewDispatchRepo(t)

	setup := func(m mocked, a args) func() {
		getCall := dispatchRepoMock.EXPECT().GetAllDispatches(a.ctx, mock.Anything).
			Maybe().Return(m.dispatches, m.getByidErr)

		return func() {
			getCall.Unset()
		}
	}

	dispatches := []e.CurrencyDispatch{{
		Id:           "id",
		Label:        "label",
		SendAt:       time.Date(1, 1, 1, 1, 1, 1, 1, time.UTC),
		TemplateName: "template",
		Details: e.CurrencyDispatchDetails{
			BaseCurrency:     "base",
			TargetCurrencies: []string{"target"}},
		CountOfSubscribers: 2,
	}}
	someErr := fmt.Errorf("some err")
	ctx := context.Background()
	arguments := args{ctx: context.Background()}
	tests := []struct {
		name    string
		args    args
		mocked  mocked
		want    []e.CurrencyDispatch
		wantErr error
	}{
		{
			name:    "failed to make transaction",
			args:    arguments,
			mocked:  mocked{getByidErr: someErr},
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
	setup := func(m mocked, a args) func() {
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
		mocked  mocked
		args    args
		wantErr error
	}{
		{
			name:    "failed to get dispatch",
			args:    a,
			mocked:  mocked{getDispatchErr: getDispatchErr},
			wantErr: getDispatchErr,
		},
		{
			name:    "user already subscribed for such dispatch",
			args:    a,
			mocked:  mocked{createUserErr: createUserErr},
			wantErr: createUserErr,
		},
		{
			name:    "failed to create subscription",
			args:    a,
			mocked:  mocked{createSubErr: createSubErr},
			wantErr: createSubErr,
		},
		{
			name: "success",
			args: a,
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
