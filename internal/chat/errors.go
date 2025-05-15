package chat

type UserIsNotAMemberError struct{}

func (e *UserIsNotAMemberError) Error() string {
	return "User is not a member of this chat"
}

type ChatDoesNotExistError struct{}

func (e *ChatDoesNotExistError) Error() string {
	return "Chat does not exist"
}
