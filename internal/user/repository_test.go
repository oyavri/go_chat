package user_test

import (
	"context"
	_ "embed"
	"fmt"
	"go_chat/internal/user"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestDbContainer struct {
	Container testcontainers.Container
	Pool      *pgxpool.Pool
}

func SetupTestDB(ctx context.Context) (*TestDbContainer, error) {
	initSQLPath, err := filepath.Abs(filepath.Join("..", "..", "db_init", "init.sql"))
	if err != nil {
		return nil, err
	}

	req := testcontainers.ContainerRequest{
		Image:        "postgres",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_DB":       "testdb",
		},
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      initSQLPath,
				ContainerFilePath: "/docker-entrypoint-initdb.d/init.sql",
				FileMode:          0o644,
			},
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

	repo := user.NewUserRepository(testDb.Pool)
	t.Run("create user successfully", func(t *testing.T) {
		username := "test_user"
		email := "test@example.org"

		testUser, err := repo.CreateUser(ctx, username, email)
		require.NoError(t, err)
		require.Equal(t, testUser.Username, username)
		require.Equal(t, testUser.Email, email)
	})

	t.Run("fail to create user with empty username", func(t *testing.T) {
		username := ""
		email := "test2@example.org"

		_, err = repo.CreateUser(ctx, username, email)
		require.Error(t, err)
		require.ErrorIs(t, err, &user.UsernameIsEmptyError{})
	})

	t.Run("fail to create user without valid email", func(t *testing.T) {
		username := "test_user3"
		email := "test"

		_, err = repo.CreateUser(ctx, username, email)
		require.Error(t, err)
	})

	t.Run("fail to create user with existing username", func(t *testing.T) {
		username := "test_user"
		email := "test4@example.org"

		_, err = repo.CreateUser(ctx, username, email)
		require.Error(t, err)
		require.ErrorIs(t, err, &user.UsernameIsTakenError{})
	})

	t.Run("create user with existing email successfully", func(t *testing.T) {
		username := "test_user5"
		email := "test@example.org"

		_, err = repo.CreateUser(ctx, username, email)
		require.NoError(t, err)
	})
}

func TestRepository_GetUserById(t *testing.T) {
	ctx := context.Background()
	testDb, err := SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDb.Terminate(ctx)

	repo := user.NewUserRepository(testDb.Pool)
	username := "test_user"
	email := "test@example.org"

	testUser, err := repo.CreateUser(ctx, username, email)
	require.NoError(t, err)

	t.Run("get existing user by id successfully", func(t *testing.T) {
		testUserFromDB, err := repo.GetUserById(ctx, testUser.Id)
		require.NoError(t, err)
		require.Equal(t, testUser.Username, testUserFromDB.Username)
	})

	t.Run("try to get non-existing user by id", func(t *testing.T) {
		notExistingUserId := uuid.New().String()
		_, err = repo.GetUserById(ctx, notExistingUserId)
		require.Error(t, err)
		require.ErrorIs(t, err, &user.UserDoesNotExistError{})
	})

}

func TestRepository_DeleteUserById(t *testing.T) {
	ctx := context.Background()
	testDb, err := SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDb.Terminate(ctx)

	repo := user.NewUserRepository(testDb.Pool)
	username := "test_user"
	email := "test@example.org"

	testUser, err := repo.CreateUser(ctx, username, email)
	require.NoError(t, err)

	t.Run("delete user by id successfully", func(t *testing.T) {
		require.Equal(t, testUser.Deleted, false)

		testUser, err = repo.DeleteUserById(ctx, testUser.Id)
		require.NoError(t, err)
		require.Equal(t, testUser.Deleted, true)
	})

	t.Run("try to delete non-existing user by id", func(t *testing.T) {
		notExistingUserId := uuid.New().String()
		_, err = repo.DeleteUserById(ctx, notExistingUserId)
		require.Error(t, err)
		require.ErrorIs(t, err, &user.UserDoesNotExistError{})
	})

	t.Run("delete already deleted user by id", func(t *testing.T) {
		require.Equal(t, testUser.Deleted, true)

		testUser, err = repo.DeleteUserById(ctx, testUser.Id)
		require.NoError(t, err)
		require.Equal(t, testUser.Deleted, true)
	})
}

func TestRepository_UpdateUserById(t *testing.T) {
	ctx := context.Background()
	testDb, err := SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDb.Terminate(ctx)

	repo := user.NewUserRepository(testDb.Pool)

	username := "test_user"
	email := "test@example.org"

	testUser, err := repo.CreateUser(ctx, username, email)
	require.NoError(t, err)

	t.Run("update username by user id successfully", func(t *testing.T) {
		newUsername := "test_user2"
		testUser, err = repo.UpdateUserById(ctx, &testUser.Id, &newUsername, nil)
		require.NoError(t, err)
		require.Equal(t, testUser.Username, newUsername)
	})

	t.Run("update email by user id successfully", func(t *testing.T) {
		newUsername := "test_user2"
		testUser, err = repo.UpdateUserById(ctx, &testUser.Id, &newUsername, nil)
		require.NoError(t, err)
		require.Equal(t, testUser.Username, newUsername)
	})

	t.Run("update both username and email by user id successfully", func(t *testing.T) {
		require.NotEqual(t, testUser.Username, username)
		require.NotEqual(t, testUser.Email, email)

		testUser, err = repo.UpdateUserById(ctx, &testUser.Id, &username, &email)
		require.NoError(t, err)
		require.Equal(t, testUser.Username, username)
		require.Equal(t, testUser.Email, email)
	})

	t.Run("update user by id without any data", func(t *testing.T) {
		testUser, err = repo.UpdateUserById(ctx, &testUser.Id, nil, nil)
		require.Error(t, err)
		require.ErrorIs(t, err, &user.NoFieldToUpdateError{})
	})

	t.Run("try to update non-existing user by id", func(t *testing.T) {
		notExistingUserId := uuid.New().String()
		newUsername := "test"

		_, err = repo.UpdateUserById(ctx, &notExistingUserId, &newUsername, nil)
		require.Error(t, err)
		require.ErrorIs(t, err, &user.UserDoesNotExistError{})
	})

	t.Run("try to update username with existing username by id", func(t *testing.T) {
		otherUsername := "test_user2"
		testUser, err := repo.CreateUser(ctx, otherUsername, email)
		require.NoError(t, err)

		testUser, err = repo.UpdateUserById(ctx, &testUser.Id, &username, nil)
		require.Error(t, err)
		require.ErrorIs(t, err, &user.UsernameIsTakenError{})
	})
}

func TestRepository_GetUserByUsername(t *testing.T) {
	ctx := context.Background()
	testDb, err := SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDb.Terminate(ctx)

	repo := user.NewUserRepository(testDb.Pool)

	username := "test_user"
	email := "test@example.org"

	testUser, err := repo.CreateUser(ctx, username, email)
	require.NoError(t, err)

	t.Run("get user by username successfully", func(t *testing.T) {
		testUserFromDB, err := repo.GetUserByUsername(ctx, testUser.Username)
		require.NoError(t, err)
		require.Equal(t, testUser.Username, testUserFromDB.Username)
		require.Equal(t, testUser.Email, testUserFromDB.Email)
	})

	t.Run("try to get non-existing user by username", func(t *testing.T) {
		notExistingUsername := "this-username-does-not-exist"
		_, err = repo.GetUserByUsername(ctx, notExistingUsername)
		require.Error(t, err)
		require.ErrorIs(t, err, &user.UserDoesNotExistError{})
	})
}

func TestRepository_DeleteUserByUsername(t *testing.T) {
	ctx := context.Background()
	testDb, err := SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDb.Terminate(ctx)

	repo := user.NewUserRepository(testDb.Pool)

	username := "test_user"
	email := "test@example.org"

	testUser, err := repo.CreateUser(ctx, username, email)
	require.NoError(t, err)

	t.Run("delete user by username successfully", func(t *testing.T) {
		require.Equal(t, testUser.Deleted, false)

		testUser, err = repo.DeleteUserByUsername(ctx, testUser.Username)
		require.NoError(t, err)
		require.Equal(t, testUser.Deleted, true)
	})

	t.Run("try to delete non-existing user by username", func(t *testing.T) {
		notExistingUsername := "this-username-does-not-exist"
		_, err = repo.DeleteUserByUsername(ctx, notExistingUsername)
		require.Error(t, err)
		require.ErrorIs(t, err, &user.UserDoesNotExistError{})
	})

	t.Run("delete already deleted user by username", func(t *testing.T) {
		require.Equal(t, testUser.Deleted, true)

		testUser, err = repo.DeleteUserByUsername(ctx, testUser.Username)
		require.NoError(t, err)
		require.Equal(t, testUser.Deleted, true)
	})
}

func TestRepository_UpdateUserByUsername(t *testing.T) {
	ctx := context.Background()
	testDb, err := SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDb.Terminate(ctx)

	repo := user.NewUserRepository(testDb.Pool)

	username := "test_user"
	email := "test@example.org"

	testUser, err := repo.CreateUser(ctx, username, email)
	require.NoError(t, err)

	t.Run("update username by username", func(t *testing.T) {
		newUsername := "test_user2"
		testUser, err = repo.UpdateUserByUsername(ctx, &testUser.Username, &newUsername, nil)
		require.NoError(t, err)
		require.Equal(t, testUser.Username, newUsername)
	})

	t.Run("update email by username", func(t *testing.T) {
		newEmail := "test2@example.org"
		testUser, err = repo.UpdateUserByUsername(ctx, &testUser.Username, nil, &newEmail)
		require.NoError(t, err)
		require.Equal(t, testUser.Email, newEmail)
	})

	t.Run("update both username and email by username", func(t *testing.T) {
		testUser, err = repo.UpdateUserByUsername(ctx, &testUser.Username, &username, &email)
		require.NoError(t, err)
		require.Equal(t, testUser.Username, username)
		require.Equal(t, testUser.Email, email)
	})

	t.Run("update user by username without any data", func(t *testing.T) {
		testUser, err = repo.UpdateUserByUsername(ctx, &testUser.Username, nil, nil)
		require.Error(t, err)
		require.ErrorIs(t, err, &user.NoFieldToUpdateError{})
	})

	t.Run("try to update non-existing user by username", func(t *testing.T) {
		notExistingUsername := "this-username-does-not-exist"
		newUsername := "test"

		_, err = repo.UpdateUserByUsername(ctx, &notExistingUsername, &newUsername, nil)
		require.Error(t, err)
		require.ErrorIs(t, err, &user.UserDoesNotExistError{})
	})

	t.Run("try to update username with existing username by username", func(t *testing.T) {
		otherUsername := "test_user2"
		testUser, err := repo.CreateUser(ctx, otherUsername, email)
		require.NoError(t, err)

		testUser, err = repo.UpdateUserByUsername(ctx, &testUser.Username, &username, nil)
		require.Error(t, err)
		require.ErrorIs(t, err, &user.UsernameIsTakenError{})
	})
}
