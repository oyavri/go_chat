package chat

type UserIsNotAMemberError struct{}

func (e *UserIsNotAMemberError) Error() string {
	return "User is not a member of this chat"
}

type ChatDoesNotExistError struct{}

func (e *ChatDoesNotExistError) Error() string {
	return "Chat does not exist"
}

type MessageContentIsEmptyError struct{}

func (e *MessageContentIsEmptyError) Error() string {
	return "Message content is empty"
}

type NoUserIdProvidedError struct{}

func (e *NoUserIdProvidedError) Error() string {
	return "No user ID provided for Chat"
}
