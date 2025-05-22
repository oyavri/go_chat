package chat

import (
	"context"
	"log/slog"
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
		slog.Error("[ChatService-CreateChat]", "Error", err)
		return Chat{}, err
	}

	return chat, nil
}

func (s *ChatService) SendMessage(ctx context.Context, req SendMessageRequest) (Message, error) {
	if ok, err := s.repo.IsMemberOfChatById(ctx, req.UserId, req.ChatId); !ok {
		slog.Error("[ChatService-CreateChat]", "Error", err)
		return Message{}, err
	}

	return s.repo.SaveMessage(ctx, req.UserId, req.ChatId, req.Content)
}

func (s *ChatService) GetMessages(ctx context.Context, chatId string, messageCount int, offset int) ([]Message, error) {
	return s.repo.GetMessages(ctx, chatId, messageCount, offset)
}
