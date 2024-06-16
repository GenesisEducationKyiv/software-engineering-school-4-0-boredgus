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

type Transaction interface {
	Rollback() error
	Commit() error
}
type WithTransaction interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (Transaction, error)
}
type Database interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}
type DatabaseWithTransaction interface {
	Database
	WithTransaction
}

type DB interface {
	IsError(err error, errCode Error) bool
	DB() Database
}

type store struct {
	database   DatabaseWithTransaction
	checkError ErrorCheckFunc
}

func NewStore(db DatabaseWithTransaction, errorF ErrorCheckFunc) *store {
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

func (s *store) DB() Database {
	return s.database
}
