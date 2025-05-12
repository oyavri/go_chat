package user

type UserDoesNotExistError struct {
	s string
}

func (e *UserDoesNotExistError) Error() string {
	return e.s
}

func (e *UserDoesNotExistError) New() error {
	return &UserDoesNotExistError{
		"User does not exist",
	}
}

type UserAlreadyExistsError struct {
	s string
}

func (e *UserAlreadyExistsError) Error() string {
	return e.s
}

func (e *UserAlreadyExistsError) New() error {
	return &UserAlreadyExistsError{
		"User already exists",
	}
}

// For later use:
type UsernameIsTakenError struct {
	s string
}

func (e *UsernameIsTakenError) Error() string {
	return e.s
}

func (e *UsernameIsTakenError) New() error {
	return &UsernameIsTakenError{
		"Username is already taken",
	}
}
