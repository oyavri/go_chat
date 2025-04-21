package user

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type GetUserRequest struct {
	UserId string `uri:"user_id"`
}

type DeleteUserRequest struct {
	UserId string `uri:"user_id"`
}

type UpdateUserRequest struct {
	UserId      string  `uri:"user_id"`
	NewUsername *string `json:"username,omitempty"`
	NewEmail    *string `json:"email,omitempty"`
}
