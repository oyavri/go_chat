package chat

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	// go:embed sql/create_chat.sql
	createChatQuery string
	// go:embed sql/save_message.sql
	saveMessageQuery string
	// go:embed sql/get_messages_by_id.sql
	getMessagesByChatIdQuery string
)

type ChatRepository struct {
	pool *pgxpool.Pool
}

func NewChatRepository(pool *pgxpool.Pool) *ChatRepository {
	return &ChatRepository{pool: pool}
}

func (r *ChatRepository) SaveChat(ctx context.Context) error {
	// Since this is not a variadic query, no need to use named args
	if _, err := r.pool.Exec(ctx, createChatQuery); err != nil {
		return err
	}

	return nil
}

func (r *ChatRepository) SaveMessage(ctx context.Context, message Message) error {
	if _, err := r.pool.Exec(ctx, saveMessageQuery, message.UserId, message.ChatId, message.Content); err != nil {
		return err
	}

	return nil
}

func (r *ChatRepository) GetMessages(ctx context.Context, chatId string, messageCount int, offset int) ([]Message, error) {
	count := messageCount

	rows, err := r.pool.Query(ctx, getMessagesByChatIdQuery, chatId, count, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var message Message
		err := rows.Scan(
			&message.Id,
			&message.UserId,
			&message.ChatId,
			&message.Content,
		)

		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
