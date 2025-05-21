package chat

import "github.com/jackc/pgx/v5/pgtype"

type Message struct {
	Id        string           `json:"id"`
	UserId    string           `json:"user_id"`
	ChatId    string           `json:"chat_id"`
	Content   string           `json:"content"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}
