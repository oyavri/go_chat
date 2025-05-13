package chat

type Message struct {
	Id        string `json:"id"`
	UserId    string `json:"user_id"`
	ChatId    string `json:"chat_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}
