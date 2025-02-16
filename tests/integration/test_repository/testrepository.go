package test_repository

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib" // Подключаем драйвер
	"github.com/stretchr/testify/require"
	testcontainers "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/pressly/goose"
)

func SetupTestDB(t *testing.T) *pgxpool.Pool {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:13",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_DB":       "testShop",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	pgContainer, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	t.Cleanup(func() { pgContainer.Terminate(context.Background()) })

	host, _ := pgContainer.Host(context.Background())
	port, _ := pgContainer.MappedPort(context.Background(), "5432")

	dsn := "postgres://postgres:postgres@" + host + ":" + port.Port() + "/testShop?sslmode=disable"
	db, err := pgxpool.New(context.Background(), dsn)
	require.NoError(t, err)

	dbMigration := stdlib.OpenDBFromPool(db)

	execMigrations(t, dbMigration)

	return db
}

func execMigrations(t *testing.T, db *sql.DB) {
	err := goose.Up(db, "migrations")

	require.NoError(t, err)
}
