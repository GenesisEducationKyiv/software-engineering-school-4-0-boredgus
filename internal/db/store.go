package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Error int

const (
	UniqueViolation Error = iota + 1
	InvalidTextRepresentation
)

type ErrorCheckFunc func(error, Error) bool

type DB interface {
	IsError(err error, errCode Error) bool
	DB() *sql.DB
}

type Store interface {
	WithTx(ctx context.Context, f func(DB) error) error
}

type store struct {
	database   *sql.DB
	checkError ErrorCheckFunc
}

func NewStore(db *sql.DB, errorF ErrorCheckFunc) Store {
	return &store{database: db, checkError: errorF}
}

func (s *store) WithTx(ctx context.Context, f func(DB) error) error {
	tx, err := s.database.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault})
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

func (s *store) IsError(err error, errCode Error) bool {
	return s.checkError(err, errCode)
}

func (s *store) DB() *sql.DB {
	return s.database
}
