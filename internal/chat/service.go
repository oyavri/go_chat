package chat

import (
	"context"
	"slices"
)

type ChatService struct {
	repo        *ChatRepository
	chatMembers map[string][]string // chatId -> []userId
}

func NewChatService(repo *ChatRepository) *ChatService {
	return &ChatService{
		repo:        repo,
		chatMembers: make(map[string][]string),
	}
}

func (s *ChatService) CreateChat(ctx context.Context, userIdList []string) (Chat, error) {
	chat, err := s.repo.SaveChat(ctx, userIdList)
	if err != nil {
		return Chat{}, err
	}

	return chat, nil
}

func (s *ChatService) SendMessage(ctx context.Context, msgReq SendMessageRequest) (Message, error) {
	if _, ok := s.chatMembers[msgReq.ChatId]; !ok {
		if isMember, err := s.repo.IsMemberOfChatById(ctx, msgReq.UserId, msgReq.ChatId); isMember {
			s.chatMembers[msgReq.ChatId] = append(s.chatMembers[msgReq.ChatId], msgReq.UserId)
		} else {
			return Message{}, err
		}
	}

	users := s.chatMembers[msgReq.ChatId]
	if !slices.Contains(users, msgReq.UserId) {
		return Message{}, &UserIsNotAMemberError{}
	}

	return s.repo.SaveMessage(ctx, msgReq)
}

func (s *ChatService) GetMessages(ctx context.Context, chatId string, messageCount int, offset int) ([]Message, error) {
	return s.repo.GetMessages(ctx, chatId, messageCount, offset)
}
