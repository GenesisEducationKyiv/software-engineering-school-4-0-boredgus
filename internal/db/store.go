package db

import (
	"context"
	"database/sql"
)

type Error int

const (
	UniqueViolation Error = iota + 1
	InvalidTextRepresentation
)

type ErrorCheckFunc func(error, Error) bool

type WithTransaction interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
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
	Database
}

type store struct {
	db         Database
	checkError ErrorCheckFunc
}

func NewStore(db Database, errorF ErrorCheckFunc) *store {
	return &store{db: db, checkError: errorF}
}

func (s *store) IsError(err error, errCode Error) bool {
	return s.checkError(err, errCode)
}

func (s *store) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return s.db.QueryRowContext(ctx, query, args...)
}

func (s *store) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return s.db.QueryContext(ctx, query, args...)
}

func (s *store) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return s.db.ExecContext(ctx, query, args...)
}
