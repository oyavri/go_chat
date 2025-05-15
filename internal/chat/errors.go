package chat

type UserIsNotAMemberError struct {
	s string
}

func (e *UserIsNotAMemberError) Error() string {
	return e.s
}

func (e *UserIsNotAMemberError) New() error {
	return &UserIsNotAMemberError{
		"User is not a member of this chat",
	}
}

type ChatDoesNotExistError struct {
	s string
}

func (e *ChatDoesNotExistError) Error() string {
	return e.s
}

func (e *ChatDoesNotExistError) New() error {
	return &ChatDoesNotExistError{
		"Chat does not exist",
	}
}
