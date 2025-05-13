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
