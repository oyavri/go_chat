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

	username := "test_user"
	email := "test@example.org"

	testUser, err := repo.CreateUser(ctx, username, email)
	require.NoError(t, err)
	require.Equal(t, testUser.Username, username)
	require.Equal(t, testUser.Email, email)
}

func TestRepository_CreateUser_WithEmptyUsername(t *testing.T) {
	ctx := context.Background()
	testDb, err := SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDb.Terminate(ctx)

	repo := user.NewUserRepository(testDb.Pool)

	username := ""
	email := "test@example.org"

	_, err = repo.CreateUser(ctx, username, email)
	require.Error(t, err)
	require.ErrorIs(t, err, &user.UsernameIsEmptyError{})
}

func TestRepository_CreateUser_WithoutValidEmail(t *testing.T) {
	ctx := context.Background()
	testDb, err := SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDb.Terminate(ctx)

	repo := user.NewUserRepository(testDb.Pool)

	username := "test_user"
	email := "test"

	_, err = repo.CreateUser(ctx, username, email)
	require.Error(t, err)
}

func TestRepository_CreateUser_WithExistingUsername(t *testing.T) {
	ctx := context.Background()
	testDb, err := SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDb.Terminate(ctx)

	repo := user.NewUserRepository(testDb.Pool)

	username := "existing_user"
	email := "test@example.org"

	_, err = repo.CreateUser(ctx, username, email)
	require.NoError(t, err)

	_, err = repo.CreateUser(ctx, username, email)
	require.Error(t, err)
	require.ErrorIs(t, err, &user.UsernameIsTakenError{})
}

func TestRepository_CreateUser_WithExistingEmail(t *testing.T) {
	ctx := context.Background()
	testDb, err := SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDb.Terminate(ctx)

	repo := user.NewUserRepository(testDb.Pool)

	username := "test_user"
	email := "test@example.org"

	_, err = repo.CreateUser(ctx, username, email)
	require.NoError(t, err)

	username = "test_user2"
	_, err = repo.CreateUser(ctx, username, email)
	require.NoError(t, err)
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

	testUserFromDB, err := repo.GetUserById(ctx, testUser.Id)
	require.NoError(t, err)
	require.Equal(t, testUser.Username, testUserFromDB.Username)
	require.Equal(t, testUser.Email, testUserFromDB.Email)
}

func TestRepository_GetUserById_NotExistingUser(t *testing.T) {
	ctx := context.Background()
	testDb, err := SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDb.Terminate(ctx)

	repo := user.NewUserRepository(testDb.Pool)

	notExistingUserId := "this-user-id-does-not-exist"
	_, err = repo.GetUserById(ctx, notExistingUserId)
	require.Error(t, err)
	require.ErrorIs(t, err, &user.UserDoesNotExistError{})
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
	require.Equal(t, testUser.Deleted, false)

	testUser, err = repo.DeleteUserById(ctx, testUser.Id)
	require.NoError(t, err)
	require.Equal(t, testUser.Deleted, true)
}

func TestRepository_DeleteUserById_NotExistingUser(t *testing.T) {
	ctx := context.Background()
	testDb, err := SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDb.Terminate(ctx)

	repo := user.NewUserRepository(testDb.Pool)

	notExistingUserId := "this-user-id-does-not-exist"
	_, err = repo.DeleteUserById(ctx, notExistingUserId)
	require.Error(t, err)
	require.ErrorIs(t, err, &user.UserDoesNotExistError{})

	notExistingUserId = uuid.New().String()
	_, err = repo.DeleteUserById(ctx, notExistingUserId)
	require.Error(t, err)
	require.ErrorIs(t, err, &user.UserDoesNotExistError{})
}

func TestRepository_DeleteUserById_AlreadyDeletedUser(t *testing.T) {
	ctx := context.Background()
	testDb, err := SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDb.Terminate(ctx)

	repo := user.NewUserRepository(testDb.Pool)

	username := "test_user"
	email := "test@example.org"

	testUser, err := repo.CreateUser(ctx, username, email)
	require.NoError(t, err)
	require.Equal(t, testUser.Deleted, false)

	testUser, err = repo.DeleteUserById(ctx, testUser.Id)
	require.NoError(t, err)
	require.Equal(t, testUser.Deleted, true)

	testUser, err = repo.DeleteUserById(ctx, testUser.Id)
	require.NoError(t, err)
	require.Equal(t, testUser.Deleted, true)

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

	newUsername := "test_user2"
	testUser, err = repo.UpdateUserById(ctx, &testUser.Id, &newUsername, nil)
	require.NoError(t, err)
	require.Equal(t, testUser.Username, newUsername)

	newEmail := "test2@example.org"
	testUser, err = repo.UpdateUserById(ctx, &testUser.Id, nil, &newEmail)
	require.NoError(t, err)
	require.Equal(t, testUser.Email, newEmail)

	testUser, err = repo.UpdateUserById(ctx, &testUser.Id, &username, &email)
	require.NoError(t, err)
	require.Equal(t, testUser.Username, username)
	require.Equal(t, testUser.Email, email)
}

func TestRepository_UpdateUserById_WithoutData(t *testing.T) {
	ctx := context.Background()
	testDb, err := SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDb.Terminate(ctx)

	repo := user.NewUserRepository(testDb.Pool)

	username := "test_user"
	email := "test@example.org"

	testUser, err := repo.CreateUser(ctx, username, email)
	require.NoError(t, err)

	testUser, err = repo.UpdateUserById(ctx, &testUser.Id, nil, nil)
	require.Error(t, err)
	require.ErrorIs(t, err, &user.NoFieldToUpdateError{})
}

func TestRepository_UpdateUserById_NotExistingUser(t *testing.T) {
	ctx := context.Background()
	testDb, err := SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDb.Terminate(ctx)

	repo := user.NewUserRepository(testDb.Pool)

	notExistingUserId := "this-user-id-does-not-exist"
	newUsername := "test"

	_, err = repo.UpdateUserById(ctx, &notExistingUserId, &newUsername, nil)
	require.Error(t, err)
	require.ErrorIs(t, err, &user.UserDoesNotExistError{})

	notExistingUserId = uuid.New().String()
	_, err = repo.UpdateUserById(ctx, &notExistingUserId, &newUsername, nil)
	require.Error(t, err)
	require.ErrorIs(t, err, &user.UserDoesNotExistError{})
}

func TestRepository_UpdateUserById_WithExistingUsername(t *testing.T) {
	ctx := context.Background()
	testDb, err := SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDb.Terminate(ctx)

	repo := user.NewUserRepository(testDb.Pool)

	username := "test_user"
	email := "test@example.org"

	_, err = repo.CreateUser(ctx, username, email)
	require.NoError(t, err)

	otherUsername := "test_user2"
	testUser, err := repo.CreateUser(ctx, otherUsername, email)
	require.NoError(t, err)

	testUser, err = repo.UpdateUserById(ctx, &testUser.Id, &username, nil)
	require.Error(t, err)
	require.ErrorIs(t, err, &user.UsernameIsTakenError{})
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

	testUserFromDB, err := repo.GetUserByUsername(ctx, testUser.Username)
	require.NoError(t, err)
	require.Equal(t, testUser.Username, testUserFromDB.Username)
	require.Equal(t, testUser.Email, testUserFromDB.Email)
}

func TestRepository_GetUserByUsername_NotExistingUser(t *testing.T) {
	ctx := context.Background()
	testDb, err := SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDb.Terminate(ctx)

	repo := user.NewUserRepository(testDb.Pool)

	notExistingUsername := "this-username-does-not-exist"
	_, err = repo.GetUserByUsername(ctx, notExistingUsername)
	require.Error(t, err)
	require.ErrorIs(t, err, &user.UserDoesNotExistError{})
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
	require.Equal(t, testUser.Deleted, false)

	testUser, err = repo.DeleteUserByUsername(ctx, testUser.Username)
	require.NoError(t, err)
	require.Equal(t, testUser.Deleted, true)

}

func TestRepository_DeleteUserByUsername_NotExistingUser(t *testing.T) {
	ctx := context.Background()
	testDb, err := SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDb.Terminate(ctx)

	repo := user.NewUserRepository(testDb.Pool)

	notExistingUsername := "this-username-does-not-exist"
	_, err = repo.DeleteUserByUsername(ctx, notExistingUsername)
	require.Error(t, err)
	require.ErrorIs(t, err, &user.UserDoesNotExistError{})
}

func TestRepository_DeleteUserByUsername_AlreadyDeletedUser(t *testing.T) {
	ctx := context.Background()
	testDb, err := SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDb.Terminate(ctx)

	repo := user.NewUserRepository(testDb.Pool)

	username := "test_user"
	email := "test@example.org"

	testUser, err := repo.CreateUser(ctx, username, email)
	require.NoError(t, err)
	require.Equal(t, testUser.Deleted, false)

	testUser, err = repo.DeleteUserByUsername(ctx, testUser.Username)
	require.NoError(t, err)
	require.Equal(t, testUser.Deleted, true)

	testUser, err = repo.DeleteUserByUsername(ctx, testUser.Username)
	require.NoError(t, err)
	require.Equal(t, testUser.Deleted, true)
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

	newUsername := "test_user2"
	testUser, err = repo.UpdateUserByUsername(ctx, &testUser.Username, &newUsername, nil)
	require.NoError(t, err)
	require.Equal(t, testUser.Username, newUsername)

	newEmail := "test2@example.org"
	testUser, err = repo.UpdateUserByUsername(ctx, &testUser.Username, nil, &newEmail)
	require.NoError(t, err)
	require.Equal(t, testUser.Email, newEmail)

	testUser, err = repo.UpdateUserByUsername(ctx, &testUser.Username, &username, &email)
	require.NoError(t, err)
	require.Equal(t, testUser.Username, username)
	require.Equal(t, testUser.Email, email)
}

func TestRepository_UpdateUserByUsername_WithoutData(t *testing.T) {
	ctx := context.Background()
	testDb, err := SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDb.Terminate(ctx)

	repo := user.NewUserRepository(testDb.Pool)

	username := "test_user"
	email := "test@example.org"

	testUser, err := repo.CreateUser(ctx, username, email)
	require.NoError(t, err)

	testUser, err = repo.UpdateUserByUsername(ctx, &testUser.Username, nil, nil)
	require.Error(t, err)
	require.ErrorIs(t, err, &user.NoFieldToUpdateError{})
}

func TestRepository_UpdateUserByUsername_NotExistingUser(t *testing.T) {
	ctx := context.Background()
	testDb, err := SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDb.Terminate(ctx)

	repo := user.NewUserRepository(testDb.Pool)

	notExistingUsername := "this-username-does-not-exist"
	newUsername := "test"

	_, err = repo.UpdateUserByUsername(ctx, &notExistingUsername, &newUsername, nil)
	require.Error(t, err)
	require.ErrorIs(t, err, &user.UserDoesNotExistError{})
}

func TestRepository_UpdateUserByUsername_WithExistingUsername(t *testing.T) {
	ctx := context.Background()
	testDb, err := SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDb.Terminate(ctx)

	repo := user.NewUserRepository(testDb.Pool)

	username := "test_user"
	email := "test@example.org"

	_, err = repo.CreateUser(ctx, username, email)
	require.NoError(t, err)

	otherUsername := "test_user2"
	testUser, err := repo.CreateUser(ctx, otherUsername, email)
	require.NoError(t, err)

	testUser, err = repo.UpdateUserByUsername(ctx, &testUser.Username, &username, nil)
	require.Error(t, err)
	require.ErrorIs(t, err, &user.UsernameIsTakenError{})
}
