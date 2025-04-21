package chat

type SendMessageRequest struct {
	ChatId  string `uri:"chat_id"`
	From    string `json:"from"`
	To      string `json:"to"`
	Content string `json:"content"`
}

type GetMessagesRequest struct {
	ChatId string `uri:"chat_id"`
}
