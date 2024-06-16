package db_mocks

import (
	context "context"
	db "subscription-api/internal/db"
)

type SimpleDB struct{}

func (d SimpleDB) IsError(err error, errCode db.Error) bool {
	return false
}
func (d SimpleDB) DB() db.Database {
	return nil
}

type Store struct {
	db db.DB
}

func (s *Store) WithTx(ctx context.Context, f func(db.DB) error) error {
	return f(s.db)
}

// NewStore creates a partitial mock of Store interface.
func NewStore() *Store {
	return &Store{db: &SimpleDB{}}
}
