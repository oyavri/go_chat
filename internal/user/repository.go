package user

import (
	"errors"
	"time"
)

type UserRepository struct {
	users map[string]User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: make(map[string]User),
	}
}

func (r *UserRepository) CreateUser(user User) error {
	if _, ok := r.users[user.Id]; ok {
		return &UserAlreadyExistsError{"User already exists"}
	}

	r.users[user.Id] = user
	return nil
}

func (r *UserRepository) GetUserById(userId string) (User, error) {
	if _, ok := r.users[userId]; !ok {
		return User{}, &UserDoesNotExistError{}
	}

	if user, ok := r.users[userId]; ok {
		return user, nil
	}

	return User{}, errors.New("unknown error when trying to get user by ID")
}

func (r *UserRepository) DeleteUserById(userId string) error {
	if _, ok := r.users[userId]; !ok {
		return &UserDoesNotExistError{}
	}

	if user, ok := r.users[userId]; ok {
		user.DeletedAt = time.Now()
		user.Deleted = true
		r.users[userId] = user
		return nil
	}

	return errors.New("unknown error when trying to delete user by ID")
}

func (r *UserRepository) UpdateUserById(userId string, newUsername *string, newEmail *string) (User, error) {
	if _, ok := r.users[userId]; !ok {
		return User{}, &UserDoesNotExistError{}
	}

	if user, ok := r.users[userId]; ok {
		if newUsername != nil {
			user.Username = *newUsername
		}
		if newEmail != nil {
			user.Email = *newEmail
		}

		user.UpdatedAt = time.Now()
		r.users[userId] = user

		user = r.users[userId]
		return user, nil
	}

	return User{}, errors.New("unknown error when trying to update user by ID")
}
