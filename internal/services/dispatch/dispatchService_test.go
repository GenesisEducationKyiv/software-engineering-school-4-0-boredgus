package ds

import (
	"context"
	"fmt"
	e "subscription-api/internal/entities"
	db_mocks "subscription-api/internal/mocks/db"
	repo_mocks "subscription-api/internal/mocks/repo"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_DispatchService_GetallDispatches(t *testing.T) {
	type mocked struct {
		txErr      error
		dispatch   e.CurrencyDispatch
		getByidErr error
	}

	storeMock := db_mocks.NewStore(t)
	dispatchRepoMock := repo_mocks.NewDispatchRepo(t)
	setup := func(m mocked) func() {
		txCall := storeMock.EXPECT().WithTx(mock.Anything, mock.Anything).Once().
			Return(m.txErr)
		getCall := dispatchRepoMock.EXPECT().GetAllDispatches(mock.Anything, mock.Anything).
			Maybe().Return(m.dispatch, m.getByidErr)

		return func() {
			txCall.Unset()
			getCall.Unset()
		}
	}

	someErr := fmt.Errorf("some err")
	tests := []struct {
		name    string
		mocked  mocked
		want    e.CurrencyDispatch
		wantErr error
	}{
		{
			name:    "failed to make transaction",
			mocked:  mocked{txErr: someErr},
			wantErr: someErr,
		},
		// TODO: test success path
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setup(tt.mocked)
			defer cleanup()

			s := &dispatchService{
				store: storeMock,
				// dispatchRepo: dispatchRepoMock,
			}
			got, err := s.GetAllDispatches(context.Background())

			assert.Equal(t, tt.want, got)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)

				return
			}
			assert.Nil(t, err)
		})
	}
}
