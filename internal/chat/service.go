package chat

import (
	"context"
)

type ChatService struct {
	repo *ChatRepository
}

func NewChatService(repo *ChatRepository) *ChatService {
	return &ChatService{
		repo: repo,
	}
}

// Is this method necessary?
func (s *ChatService) CreateChat(ctx context.Context) (Chat, error) {
	err := s.repo.SaveChat(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *ChatService) SendMessage(ctx context.Context, request SendMessageRequest) (Message, error) {
	// add cache for existance of the chat
	message := Message{
		UserId:  request.UserId,
		ChatId:  request.ChatId,
		Content: request.Content,
	}

	err := s.repo.SaveMessage(ctx, message)
	return message, err
}

func (s *ChatService) GetMessages(ctx context.Context, chatId string, messageCount int, offset int) ([]Message, error) {
	return s.repo.GetMessages(ctx, chatId, messageCount, offset)
}
