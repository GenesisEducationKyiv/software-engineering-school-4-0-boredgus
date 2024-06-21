package tests

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	Database string = "subs"
	Username string = "postgres"
	Password string = "pass"
)

type PostgresContainer struct {
	*postgres.PostgresContainer
	ConnectionString string
}

func (c *PostgresContainer) ExecuteSQLFiles(ctx context.Context, filenames ...string) error {
	var buffer bytes.Buffer
	for _, filename := range filenames {
		filenameWithExtension := fmt.Sprintf("/%s.sql", filename)

		err := c.PostgresContainer.CopyFileToContainer(
			ctx,
			fmt.Sprintf("./testdata%s", filenameWithExtension),
			filenameWithExtension, 64)
		if err != nil {
			return err
		}

		_, reader, err := c.PostgresContainer.Exec(ctx,
			[]string{"psql", "-U", Username, "-d", Database, "-a", "-f", filenameWithExtension})
		if err != nil {
			return err
		}

		buffer.Reset()
		if _, err = buffer.ReadFrom(reader); err != nil {
			return err
		}
		fmt.Printf("\n> EXECUTED FILE %s:\n\n%s\n\n", filenameWithExtension, buffer.String())
	}

	return nil
}

func CreatePostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16.2-bullseye"),
		postgres.WithDatabase(Database),
		postgres.WithUsername(Username),
		postgres.WithPassword(Password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second),
		))
	if err != nil {
		return nil, err
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable", "search_path="+Database)
	if err != nil {
		return nil, err
	}

	return &PostgresContainer{
		PostgresContainer: pgContainer,
		ConnectionString:  connStr,
	}, nil
}
