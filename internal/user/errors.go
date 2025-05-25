package user

type UserDoesNotExistError struct{}

func (e *UserDoesNotExistError) Error() string {
	return "User does not exist"
}

type UsernameIsTakenError struct{}

func (e *UsernameIsTakenError) Error() string {
	return "Username is already taken"
}

type UsernameIsEmptyError struct{}

func (e *UsernameIsEmptyError) Error() string {
	return "Username cannot be empty"
}

type NoFieldToUpdateError struct{}

func (e *NoFieldToUpdateError) Error() string {
	return "No field to update"
}
