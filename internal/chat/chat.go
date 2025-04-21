package chat

type Chat struct {
	Id       string    `json:"id"`
	Messages []Message `json:"messages"`
	// Members  []User    `json:"members"`
}
