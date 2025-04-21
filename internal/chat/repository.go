package chat

type ChatRepository struct {
	chats map[string]Chat
}

func NewChatRepository() *ChatRepository {
	return &ChatRepository{
		chats: make(map[string]Chat),
	}
}

func (r *ChatRepository) SaveChat(chat Chat) error {
	r.chats[chat.Id] = chat
	return nil
}

func (r *ChatRepository) SaveMessage(chatId string, message Message) error {
	chat := r.chats[chatId]
	chat.Messages = append(chat.Messages, message)
	r.chats[chatId] = chat
	return nil
}

func (r *ChatRepository) GetMessages(chatId string) ([]Message, error) {
	chat := r.chats[chatId]
	return chat.Messages, nil
}
