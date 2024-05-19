package sql

import (
	"database/sql"
	"embed"
	"fmt"
	cfg "subscription-api/config"

	"github.com/pressly/goose/v3"
)

func initMigrations(dialect string, db *sql.DB) {
	if err := goose.SetDialect(dialect); err != nil {
		cfg.Log().Fatalf("failed to set %v dialect: %v", dialect, err)
	}
	if version, err := goose.EnsureDBVersion(db); err == nil {
		cfg.Log().Infof("version of %s migrations: %v", dialect, version)
	}
	if err := goose.Up(db, dialect); err != nil {
		cfg.Log().Fatalf("failed to make %v migrations up: %v", dialect, err)
	}
}

//go:embed postgres/*.sql
var potgresqlMigrations embed.FS

func PostgeSQLMigrationsUp(scheme string) func(db *sql.DB) {
	return func(db *sql.DB) {
		goose.SetBaseFS(potgresqlMigrations)
		goose.SetTableName(fmt.Sprintf("%s.goose_db_version", scheme))
		initMigrations(string(goose.DialectPostgres), db)
	}
}
