package sql

import (
	"database/sql"
	"embed"
	"fmt"
	"subscription-api/pkg/utils"

	"github.com/pressly/goose/v3"
)

func initMigrations(dialect string, db *sql.DB) {
	if err := goose.SetDialect(dialect); err != nil {
		utils.PanicOnError(err, fmt.Sprintf("failed to set %v dialect", dialect))
	}
	if err := goose.Up(db, dialect); err != nil {
		utils.PanicOnError(err, fmt.Sprintf("failed to make %v migrations up", dialect))
	}
}

//go:embed postgres/*.sql
var potgresqlMigrations embed.FS

func PostgeSQLMigrationsUp(scheme string, l goose.Logger) func(db *sql.DB) {
	return func(db *sql.DB) {
		goose.SetLogger(l)
		goose.SetBaseFS(potgresqlMigrations)
		goose.SetTableName(fmt.Sprintf("%s.goose_db_version", scheme))
		initMigrations(string(goose.DialectPostgres), db)
	}
}
