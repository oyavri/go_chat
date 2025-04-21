package chat

import (
	"github.com/google/uuid"
)

type ChatService struct {
	repo *ChatRepository
}

func NewChatService(repo *ChatRepository) *ChatService {
	return &ChatService{
		repo: repo,
	}
}

func (s *ChatService) CreateChat() (Chat, error) {
	chat := Chat{
		Id:       uuid.New().String(),
		Messages: []Message{},
	}

	err := s.repo.SaveChat(chat)
	return chat, err
}

func (s *ChatService) SendMessage(request SendMessageRequest) (Message, error) {
	message := Message{
		Id:      uuid.New().String(),
		From:    request.From,
		To:      request.To,
		Content: request.Content,
	}

	err := s.repo.SaveMessage(request.ChatId, message)
	return message, err
}

func (s *ChatService) GetMessages(chatId string) ([]Message, error) {
	return s.repo.GetMessages(chatId)
}
