package chat

type Chat struct {
	Id      string   `json:"id"`
	Members []string `json:"members"`
}
