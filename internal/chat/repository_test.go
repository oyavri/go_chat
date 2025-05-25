package chat_test

import (
	"context"
	"fmt"
	"go_chat/internal/chat"
	"go_chat/internal/user"
	"path/filepath"
	"testing"

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

func TestRepository_CreateChat_WithSingleUser(t *testing.T) {
	ctx := context.Background()
	testDb, err := SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDb.Terminate(ctx)

	chatRepo := chat.NewChatRepository(testDb.Pool)
	userRepo := user.NewUserRepository(testDb.Pool)

	username := "test_user"
	email := "test@example.org"
	testUser, err := userRepo.CreateUser(ctx, username, email)
	require.NoError(t, err)

	_, err = chatRepo.SaveChat(ctx, []string{testUser.Id})
	require.NoError(t, err)
}

func TestRepository_CreateChat_WithMoreThanOneUser(t *testing.T) {
	ctx := context.Background()
	testDb, err := SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDb.Terminate(ctx)

	chatRepo := chat.NewChatRepository(testDb.Pool)
	userRepo := user.NewUserRepository(testDb.Pool)

	username := "test_user"
	email := "test@example.org"
	firstUser, err := userRepo.CreateUser(ctx, username, email)
	require.NoError(t, err)

	username = "test_user2"
	secondUser, err := userRepo.CreateUser(ctx, username, email)
	require.NoError(t, err)

	_, err = chatRepo.SaveChat(ctx, []string{firstUser.Id, secondUser.Id})
	require.NoError(t, err)
}

func TestRepository_CreateChat_WithoutUser(t *testing.T) {
	ctx := context.Background()
	testDb, err := SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDb.Terminate(ctx)

	chatRepo := chat.NewChatRepository(testDb.Pool)
	_, err = chatRepo.SaveChat(ctx, []string{})
	require.Error(t, err)
	require.ErrorIs(t, err, &chat.NoUserIdProvidedError{})
}

func TestRepository_SaveMessage(t *testing.T) {
	ctx := context.Background()
	testDb, err := SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDb.Terminate(ctx)

	chatRepo := chat.NewChatRepository(testDb.Pool)
	userRepo := user.NewUserRepository(testDb.Pool)

	username := "test_user"
	email := "test@example.org"
	testUser, err := userRepo.CreateUser(ctx, username, email)
	require.NoError(t, err)

	chat, err := chatRepo.SaveChat(ctx, []string{testUser.Id})
	require.NoError(t, err)

	testMessage := "This is a test message"
	message, err := chatRepo.SaveMessage(ctx, testUser.Id, chat.Id, testMessage)

	require.NoError(t, err)
	require.Equal(t, message.ChatId, chat.Id)
	require.Equal(t, message.UserId, testUser.Id)
}

func TestRepository_SaveMessage_EmptyMessage(t *testing.T) {
	ctx := context.Background()
	testDb, err := SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDb.Terminate(ctx)

	chatRepo := chat.NewChatRepository(testDb.Pool)
	userRepo := user.NewUserRepository(testDb.Pool)

	username := "test_user"
	email := "test@example.org"
	testUser, err := userRepo.CreateUser(ctx, username, email)
	require.NoError(t, err)

	c, err := chatRepo.SaveChat(ctx, []string{testUser.Id})
	require.NoError(t, err)

	testMessage := ""
	_, err = chatRepo.SaveMessage(ctx, testUser.Id, c.Id, testMessage)

	require.Error(t, err)
	require.ErrorIs(t, err, &chat.MessageContentIsEmptyError{})
}

func TestRepository_GetMessages(t *testing.T) {
	ctx := context.Background()
	testDb, err := SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDb.Terminate(ctx)

	chatRepo := chat.NewChatRepository(testDb.Pool)
	userRepo := user.NewUserRepository(testDb.Pool)

	username := "test_user"
	email := "test@example.org"
	testUser, err := userRepo.CreateUser(ctx, username, email)
	require.NoError(t, err)

	c, err := chatRepo.SaveChat(ctx, []string{testUser.Id})
	require.NoError(t, err)

	testMessage := "This is a test message"
	message, err := chatRepo.SaveMessage(ctx, testUser.Id, c.Id, testMessage)

	require.NoError(t, err)
	require.Equal(t, message.ChatId, c.Id)
	require.Equal(t, message.UserId, testUser.Id)

	msgCount := 30
	offset := 0

	messages, err := chatRepo.GetMessages(ctx, c.Id, msgCount, offset)
	require.NoError(t, err)
	require.Equal(t, messages[0].Content, message.Content)
}

func TestRepository_GetMessages_WithOffset(t *testing.T) {
	ctx := context.Background()
	testDb, err := SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDb.Terminate(ctx)

	chatRepo := chat.NewChatRepository(testDb.Pool)
	userRepo := user.NewUserRepository(testDb.Pool)

	username := "test_user"
	email := "test@example.org"
	testUser, err := userRepo.CreateUser(ctx, username, email)
	require.NoError(t, err)

	c, err := chatRepo.SaveChat(ctx, []string{testUser.Id})
	require.NoError(t, err)

	testMessage := "This is the first test message"
	_, err = chatRepo.SaveMessage(ctx, testUser.Id, c.Id, testMessage)
	require.NoError(t, err)

	testMessage = "This is the second test message"
	message, err := chatRepo.SaveMessage(ctx, testUser.Id, c.Id, testMessage)
	require.NoError(t, err)
	require.Equal(t, message.ChatId, c.Id)
	require.Equal(t, message.UserId, testUser.Id)

	msgCount := 1
	offset := 1

	messages, err := chatRepo.GetMessages(ctx, c.Id, msgCount, offset)
	require.NoError(t, err)
	require.Equal(t, messages[0].Content, message.Content)
}

func TestRepository_GetMessages_WithoutMessage(t *testing.T) {
	ctx := context.Background()
	testDb, err := SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDb.Terminate(ctx)

	chatRepo := chat.NewChatRepository(testDb.Pool)
	userRepo := user.NewUserRepository(testDb.Pool)

	username := "test_user"
	email := "test@example.org"
	testUser, err := userRepo.CreateUser(ctx, username, email)
	require.NoError(t, err)

	c, err := chatRepo.SaveChat(ctx, []string{testUser.Id})
	require.NoError(t, err)

	msgCount := 30
	offset := 0

	_, err = chatRepo.GetMessages(ctx, c.Id, msgCount, offset)
	require.NoError(t, err)
}

func TestRepository_IsMemberOfChatById(t *testing.T) {
	ctx := context.Background()
	testDb, err := SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDb.Terminate(ctx)

	chatRepo := chat.NewChatRepository(testDb.Pool)
	userRepo := user.NewUserRepository(testDb.Pool)

	username := "test_user"
	email := "test@example.org"
	testUser, err := userRepo.CreateUser(ctx, username, email)
	require.NoError(t, err)

	c, err := chatRepo.SaveChat(ctx, []string{testUser.Id})
	require.NoError(t, err)

	isMember, err := chatRepo.IsMemberOfChatById(ctx, testUser.Id, c.Id)
	require.NoError(t, err)
	require.Equal(t, isMember, true)

	notExistingUserId := "this-user-id-does-not-exist"

	isMember, err = chatRepo.IsMemberOfChatById(ctx, notExistingUserId, c.Id)
	require.Error(t, err)
	require.ErrorIs(t, err, &chat.UserIsNotAMemberError{})
	require.Equal(t, isMember, false)
}

func TestRepository_GetChatById(t *testing.T) {
	ctx := context.Background()
	testDb, err := SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDb.Terminate(ctx)

	chatRepo := chat.NewChatRepository(testDb.Pool)
	userRepo := user.NewUserRepository(testDb.Pool)

	username := "test_user"
	email := "test@example.org"
	testUser, err := userRepo.CreateUser(ctx, username, email)
	require.NoError(t, err)

	c, err := chatRepo.SaveChat(ctx, []string{testUser.Id})
	require.NoError(t, err)

	isChatExist, err := chatRepo.GetChatById(ctx, c.Id)
	require.NoError(t, err)
	require.Equal(t, isChatExist, true)
}

func TestRepository_GetChatById_NotExistingUser(t *testing.T) {
	ctx := context.Background()
	testDb, err := SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDb.Terminate(ctx)

	chatRepo := chat.NewChatRepository(testDb.Pool)

	notExistingChatId := "this-chat-id-does-not-exist"

	isChatExist, err := chatRepo.GetChatById(ctx, notExistingChatId)
	require.Error(t, err)
	require.Equal(t, isChatExist, false)
}
