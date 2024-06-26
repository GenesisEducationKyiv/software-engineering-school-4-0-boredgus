package sql

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
)

func initMigrations(dialect string, db *sql.DB) error {
	if err := goose.SetDialect(dialect); err != nil {
		return fmt.Errorf("failed to set %v dialect: %w", dialect, err)
	}
	if err := goose.Up(db, dialect); err != nil {
		return fmt.Errorf("failed to make %v migrations up: %w", dialect, err)
	}

	return nil
}

//go:embed postgres/*.sql
var potgresqlMigrations embed.FS

func PostgeSQLMigrationsUp(l goose.Logger) func(db *sql.DB) error {
	return func(db *sql.DB) error {
		if l != nil {
			goose.SetLogger(l)
		}
		goose.SetBaseFS(potgresqlMigrations)
		goose.SetTableName("public.goose_db_version")

		return initMigrations(string(goose.DialectPostgres), db)
	}
}
