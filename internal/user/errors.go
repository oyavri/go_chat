package user

type UserDoesNotExistError struct{}

func (e *UserDoesNotExistError) Error() string {
	return "User does not exist"
}

type UserAlreadyExistsError struct{}

func (e *UserAlreadyExistsError) Error() string {
	return "User already exists"
}

type UsernameIsTakenError struct{}

func (e *UsernameIsTakenError) Error() string {
	return "Username is already taken"
}
