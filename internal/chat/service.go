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

func (s *ChatService) CreateChat(ctx context.Context, chatReq CreateChatRequest) (Chat, error) {
	chat, err := s.repo.SaveChat(ctx, chatReq.Members)
	if err != nil {
		return Chat{}, err
	}

	return chat, nil
}

func (s *ChatService) SendMessage(ctx context.Context, msgReq SendMessageRequest) (Message, error) {
	if ok, err := s.repo.IsMemberOfChatById(ctx, msgReq.UserId, msgReq.ChatId); !ok {
		return Message{}, err
	}

	return s.repo.SaveMessage(ctx, msgReq)
}

func (s *ChatService) GetMessages(ctx context.Context, chatId string, messageCount int, offset int) ([]Message, error) {
	return s.repo.GetMessages(ctx, chatId, messageCount, offset)
}
