package tests

import (
	"bytes"
	"context"
	"fmt"
	"path"
	"runtime"
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
	ConnectionString  string
	testDataDirectory string
}

func (c *PostgresContainer) ExecuteSQLFiles(ctx context.Context, filenames ...string) error {
	var buffer bytes.Buffer
	for _, filename := range filenames {
		filenameWithExtension := fmt.Sprintf("%s.sql", filename)
		localFilePath := path.Join(c.testDataDirectory, filenameWithExtension)
		targetFilePath := path.Join("data", filenameWithExtension)

		err := c.PostgresContainer.CopyFileToContainer(ctx, localFilePath, targetFilePath, 64)
		if err != nil {
			return err
		}

		_, reader, err := c.PostgresContainer.Exec(ctx,
			[]string{"psql", "-U", Username, "-d", Database, "-a", "-f", localFilePath})
		if err != nil {
			return err
		}

		buffer.Reset()
		if _, err = buffer.ReadFrom(reader); err != nil {
			return err
		}
		fmt.Printf("\n> EXECUTED FILE %s:\n\n%s\n\n", localFilePath, buffer.String())
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

	_, filename, _, _ := runtime.Caller(0)

	return &PostgresContainer{
		PostgresContainer: pgContainer,
		ConnectionString:  connStr,
		testDataDirectory: path.Join(path.Dir(filename), "testdata"),
	}, nil
}
