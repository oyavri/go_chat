package chat

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ChatRepository struct {
	pool *pgxpool.Pool
}

func NewChatRepository(pool *pgxpool.Pool) *ChatRepository {
	return &ChatRepository{pool: pool}
}

func (r *ChatRepository) SaveChat(ctx context.Context) error {
	// Since this is not a variadic query, no need to use named args
	query := `INSERT INTO chat (id) 
			  VALUES (DEFAULT)`

	if _, err := r.pool.Exec(ctx, query); err != nil {
		return err
	}

	return nil
}

func (r *ChatRepository) SaveMessage(ctx context.Context, message Message) error {
	query := `INSERT INTO chat_message (user_id, chat_id, content) 
			  VALUES ($1, $2, $3)`

	if _, err := r.pool.Exec(ctx, query, message.UserId, message.ChatId, message.Content); err != nil {
		return err
	}

	return nil
}

func (r *ChatRepository) GetMessages(ctx context.Context, chatId string, messageCount int, offset int) ([]Message, error) {
	query := `SELECT * FROM chat_message 
			  JOIN chat ON chat_message.chat_id = chat.id 
			  WHERE chat.id = $1 
			  ORDER BY chat_message.created_at 
			  LIMIT $2 
			  OFFSET $3`

	count := messageCount
	// Hard-coded limit for now.
	if count > 30 {
		count = 30
	}

	rows, err := r.pool.Query(ctx, query, chatId, count, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var message Message
		err := rows.Scan(&message.Id, &message.UserId, &message.ChatId, &message.Content)
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
