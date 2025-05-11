package chat

type SendMessageRequest struct {
	ChatId  string `json:"chat_id"`
	UserId  string `json:"user_id"`
	Content string `json:"content"`
}

type GetMessagesRequest struct {
	ChatId       string `uri:"chat_id"`
	MessageCount int    `uri:"message_count"`
	Offset       int    `uri:"offset"`
}
