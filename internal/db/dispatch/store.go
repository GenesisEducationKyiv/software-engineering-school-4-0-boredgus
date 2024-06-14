package ds

import (
	"context"
	"database/sql"
	"fmt"
	"subscription-api/internal/db"
	"subscription-api/internal/entities"
)

type DispatchData struct {
	entities.Dispatch[entities.CurrencyDispatchDetails]
	CountOfSubscribers int
}

type SubscribeForParams struct {
	Email, Dispatch string
}

type DispatchQueries interface {
	GetDispatch(ctx context.Context, id string) (DispatchData, error)
	SubscribeFor(context.Context, SubscribeForParams) error
	CreateSubscriber(ctx context.Context, email string) error
}

type DispatchStore interface {
	DispatchQueries
	db.Store[DispatchQueries]
}

type currencyDispatchStore struct {
	DispatchQueries
	store db.Store[any]
}

func NewCurrencyDispatchStore(store db.Store[any]) DispatchStore {
	return &currencyDispatchStore{store: store}
}
func (s *currencyDispatchStore) WithTx(ctx context.Context, f func(q DispatchQueries) error) error {
	tx, err := s.store.DB().BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	if err = f(s); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("failed to rollback transaction: %w: %w", err, rbErr)
		}

		return err
	}

	return tx.Commit()
}

func (s *currencyDispatchStore) DB() *sql.DB {
	return s.store.DB()
}
func (s *currencyDispatchStore) IsError(err error, errCode db.Error) bool {
	return s.store.IsError(err, errCode)
}
