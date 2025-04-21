package user

import (
	"time"

	"github.com/google/uuid"
)

type UserService struct {
	repo *UserRepository
}

func NewUserService(repo *UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) CreateUser(req CreateUserRequest) (User, error) {
	user := User{
		Id:        uuid.New().String(),
		Username:  req.Username,
		Email:     req.Email,
		CreatedAt: time.Now(),
		Deleted:   false,
	}

	err := s.repo.CreateUser(user)
	return user, err
}

func (s *UserService) GetUserById(userId string) (User, error) {
	return s.repo.GetUserById(userId)
}

func (s *UserService) DeleteUserById(userId string) error {
	return s.repo.DeleteUserById(userId)
}

func (s *UserService) UpdateUserById(userId string, newUsername *string, newEmail *string) (User, error) {
	return s.repo.UpdateUserById(userId, newUsername, newEmail)
}
