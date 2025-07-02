package chat_test

import (
	"context"
	"fmt"
	"go_chat/internal/chat"
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

func TestRepository_CreateChat(t *testing.T) {
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

	t.Run("create chat with single user", func(t *testing.T) {
		_, err = chatRepo.SaveChat(ctx, []string{testUser.Id})
		require.NoError(t, err)
	})

	t.Run("create chat with multiple users", func(t *testing.T) {
		username := "test_user2"
		otherUser, err := userRepo.CreateUser(ctx, username, email)
		require.NoError(t, err)

		_, err = chatRepo.SaveChat(ctx, []string{testUser.Id, otherUser.Id})
		require.NoError(t, err)
	})

	t.Run("create chat with no user", func(t *testing.T) {
		_, err = chatRepo.SaveChat(ctx, []string{})
		require.Error(t, err)
		require.ErrorIs(t, err, &chat.NoUserIdProvidedError{})
	})
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

	c, err := chatRepo.SaveChat(ctx, []string{testUser.Id})
	require.NoError(t, err)

	t.Run("save message", func(t *testing.T) {
		testMessage := "This is a test message"
		message, err := chatRepo.SaveMessage(ctx, testUser.Id, c.Id, testMessage)

		require.NoError(t, err)
		require.Equal(t, message.ChatId, c.Id)
		require.Equal(t, message.UserId, testUser.Id)
	})

	t.Run("save empty message", func(t *testing.T) {
		testMessage := ""
		_, err = chatRepo.SaveMessage(ctx, testUser.Id, c.Id, testMessage)

		require.Error(t, err)
		require.ErrorIs(t, err, &chat.MessageContentIsEmptyError{})
	})
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

	testMessage := "This is the first test message"
	_, err = chatRepo.SaveMessage(ctx, testUser.Id, c.Id, testMessage)
	require.NoError(t, err)

	testMessage = "This is the second test message"
	message, err := chatRepo.SaveMessage(ctx, testUser.Id, c.Id, testMessage)
	require.NoError(t, err)
	require.Equal(t, message.ChatId, c.Id)
	require.Equal(t, message.UserId, testUser.Id)

	t.Run("get messages", func(t *testing.T) {
		msgCount := 30
		offset := 0

		messages, err := chatRepo.GetMessages(ctx, c.Id, msgCount, offset)
		require.NoError(t, err)
		require.Equal(t, messages[0].Content, message.Content)
	})

	t.Run("get messages with offset", func(t *testing.T) {
		msgCount := 1
		offset := 1

		messages, err := chatRepo.GetMessages(ctx, c.Id, msgCount, offset)
		require.NoError(t, err)
		require.Equal(t, messages[0].Content, message.Content)
	})

	t.Run("get messages without message", func(t *testing.T) {
		c, err := chatRepo.SaveChat(ctx, []string{testUser.Id})
		require.NoError(t, err)

		msgCount := 30
		offset := 0

		_, err = chatRepo.GetMessages(ctx, c.Id, msgCount, offset)
		require.NoError(t, err)
	})
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

	t.Run("is member of chat by id", func(t *testing.T) {
		isMember, err := chatRepo.IsMemberOfChatById(ctx, testUser.Id, c.Id)
		require.NoError(t, err)
		require.Equal(t, isMember, true)

		notExistingUserId := "this-user-id-does-not-exist"

		isMember, err = chatRepo.IsMemberOfChatById(ctx, notExistingUserId, c.Id)
		require.Error(t, err)
		require.ErrorIs(t, err, &chat.UserIsNotAMemberError{})
		require.Equal(t, isMember, false)
	})
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

	t.Run("get chat by id", func(t *testing.T) {
		isChatExist, err := chatRepo.GetChatById(ctx, c.Id)
		require.NoError(t, err)
		require.Equal(t, isChatExist, true)
	})

	t.Run("get chat by non-existing id", func(t *testing.T) {
		notExistingChatId := uuid.New().String()

		isChatExist, err := chatRepo.GetChatById(ctx, notExistingChatId)
		require.Error(t, err)
		require.Equal(t, isChatExist, false)
	})
}
