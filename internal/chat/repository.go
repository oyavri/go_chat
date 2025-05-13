package chat

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	// go:embed sql/create_chat.sql
	createChatQuery string
	// go:embed sql/add_chat_member_by_id.sql
	addChatMemberByIdQuery string
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

func (r *ChatRepository) SaveChat(ctx context.Context, userIdList []string) (Chat, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return Chat{}, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	var chatId string
	err = tx.QueryRow(ctx, createChatQuery).
		Scan(&chatId)
	if err != nil {
		return Chat{}, err
	}

	var insertedUserIdList []string
	var addedUserId string

	for _, userId := range userIdList {
		err := tx.QueryRow(ctx, addChatMemberByIdQuery, chatId, userId).
			Scan(&addedUserId)
		if err != nil {
			return Chat{}, err
		}

		insertedUserIdList = append(insertedUserIdList, addedUserId)
	}

	return Chat{
		Id:      chatId,
		Members: insertedUserIdList,
	}, nil
}

func (r *ChatRepository) SaveMessage(ctx context.Context, msgReq SendMessageRequest) (Message, error) {
	var message Message

	err := r.pool.QueryRow(ctx, saveMessageQuery, msgReq.UserId, msgReq.ChatId, msgReq.Content).
		Scan(
			&message.Id,
			&message.UserId,
			&message.ChatId,
			&message.Content,
			&message.CreatedAt,
		)
	if err != nil {
		return Message{}, err
	}

	return message, nil
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
			&message.CreatedAt,
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
