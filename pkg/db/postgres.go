package db

import (
	"database/sql"
	"fmt"
	"subscription-api/internal/db"
	"sync"

	"github.com/lib/pq"
)

var postgresOnce sync.Once
var postgresDB *sql.DB

func NewPostrgreSQL(dsn string, handlers ...func(db *sql.DB)) (*sql.DB, error) {
	postgresOnce.Do(func() {
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			panic(fmt.Errorf("failed to connect to postgres db: %w", err))
			// return nil, err
		}
		postgresDB = db
	})
	for _, handler := range handlers {
		handler(postgresDB)
	}
	return postgresDB, nil
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
	// fmt.Println("> IsPqError: \n", e.Code, strings.Join([]string{e.Message, e.Detail, e.Hint, e.Line}, ";"))
	if !ok || e.Code != pqErrors[errCode] {
		return false
	}
	return true
}
