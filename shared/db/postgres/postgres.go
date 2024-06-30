package postgres

import (
	"database/sql"
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/shared/db"
	"github.com/lib/pq"
)

func NewPostrgreSQL(dsn string, handlers ...func(db *sql.DB) error) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(fmt.Errorf("failed to connect to postgres db: %w", err))
	}

	for _, handler := range handlers {
		if err = handler(db); err != nil {
			return db, nil
		}
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
	pqErr, ok := err.(*pq.Error)
	if !ok || pqErr.Code != pqErrors[errCode] {
		return false
	}

	return true
}
