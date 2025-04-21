package user

import (
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
	if user, ok := r.users[userId]; ok {
		return user, nil
	}

	return User{}, &UserDoesNotExistError{"User does not exist"}
}

func (r *UserRepository) DeleteUserById(userId string) error {
	if user, ok := r.users[userId]; ok {
		user.DeletedAt = time.Now()
		user.Deleted = true
		r.users[userId] = user
		return nil
	}

	return &UserDoesNotExistError{"User does not exist"}
}

// Not sure how I should change the user
func (r *UserRepository) UpdateUserById(userId string, newUsername *string, newEmail *string) (User, error) {
	// Would it be better to carry the code inside the if clause
	// to outside to prevent if-hadouken (IDK if it occurs in Go)
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

	return User{}, &UserDoesNotExistError{"User does not exist"}
}
