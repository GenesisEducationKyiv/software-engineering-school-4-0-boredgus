package db

import (
	"database/sql"
	"fmt"
	"subscription-api/internal/db"

	"github.com/lib/pq"
)

func NewPostrgreSQL(dsn string, handlers ...func(db *sql.DB)) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(fmt.Errorf("failed to connect to postgres db: %w", err))
	}

	for _, handler := range handlers {
		handler(db)
	}

	return db, nil
}

const (
	UniqueViolationError                          = pq.ErrorCode("23505") // 'unique_violation'
	SchemaAndDataStatementMixingNotSupportedError = pq.ErrorCode("25007") // 'schema_and_data_statement_mixing_not_supported'
	InvalidTextRepresentation                     = pq.ErrorCode("22P02") // 'invalid_text_representation'
)

var pqErrors = map[db.Error]pq.ErrorCode{
	db.UniqueViolation:           UniqueViolationError,
	db.InvalidTextRepresentation: InvalidTextRepresentation,
}

func IsPqError(err error, errCode db.Error) bool {
	e, ok := err.(*pq.Error)
	if !ok || e.Code != pqErrors[errCode] {
		return false
	}

	return true
}
