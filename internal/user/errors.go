package user

type UserDoesNotExistError struct {
	s string
}

func (e *UserDoesNotExistError) Error() string {
	return e.s
}

func (e *UserDoesNotExistError) New(message string) error {
	return &UserDoesNotExistError{
		message,
	}
}

// Should HTTP codes of the errors exist in this domain?
type UserAlreadyExistsError struct {
	s string
}

func (e *UserAlreadyExistsError) Error() string {
	return e.s
}

// The name is already descriptive, should the message be variable?
func (e *UserAlreadyExistsError) New(message string) error {
	return &UserAlreadyExistsError{
		message,
	}
}
