package user

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type GetUserRequest struct {
	UserId string `params:"user_id"`
}

type DeleteUserRequest struct {
	UserId string `params:"user_id"`
}

type UpdateUserRequest struct {
	UserId      string  `params:"user_id"`
	NewUsername *string `json:"username"`
	NewEmail    *string `json:"email"`
}
