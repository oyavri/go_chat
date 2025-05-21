package user

import (
	"context"
)

type UserService struct {
	repo *UserRepository
}

func NewUserService(repo *UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) (User, error) {
	user, err := s.repo.CreateUser(ctx, req.Username, req.Email)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (s *UserService) GetUserById(ctx context.Context, userId string) (User, error) {
	return s.repo.GetUserById(ctx, userId)
}

func (s *UserService) DeleteUserById(ctx context.Context, userId string) (User, error) {
	return s.repo.DeleteUserById(ctx, userId)
}

func (s *UserService) UpdateUserById(ctx context.Context, userId *string, newUsername *string, newEmail *string) (User, error) {
	return s.repo.UpdateUserById(ctx, userId, newUsername, newEmail)
}
