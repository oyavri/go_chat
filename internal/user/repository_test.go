package user_test

import (
	"context"
	_ "embed"
	"fmt"
	"go_chat/internal/user"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	// go:embed db_init/init.sql
	initDbScript string
)

type TestDbContainer struct {
	Container testcontainers.Container
	Pool      *pgxpool.Pool
}

func SetupTestDB(ctx context.Context) (*TestDbContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("postgres://testuser:testpass@%s:%s/testdb?sslmode=disable", host, port.Port())

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	return &TestDbContainer{
		Container: container,
		Pool:      pool,
	}, nil
}

func (c *TestDbContainer) Terminate(ctx context.Context) error {
	c.Pool.Close()
	return c.Container.Terminate(ctx)
}

func TestRepository_CreateUser(t *testing.T) {
	ctx := context.Background()
	testDb, err := SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDb.Terminate(ctx)

	_, err = testDb.Pool.Exec(ctx, initDbScript)
	require.NoError(t, err)

	repo := user.NewUserRepository(testDb.Pool)

	username := "test_user"
	email := "test@example.org"

	testUser, err := repo.CreateUser(ctx, username, email)
	require.NoError(t, err)
	require.Equal(t, testUser.Username, username)
	require.Equal(t, testUser.Email, email)
}

func TestRepository_GetUserById(t *testing.T) {

}

func TestRepository_DeleteUserById(t *testing.T) {

}

func TestRepository_UpdateUserById(t *testing.T) {

}

func TestRepository_GetUserByUsername(t *testing.T) {

}

func TestRepository_DeleteUserByUsername(t *testing.T) {

}

func TestRepository_UpdateUserByUsername(t *testing.T) {

}
