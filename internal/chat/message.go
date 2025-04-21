package chat

type Message struct {
	Id      string `json:"id"`
	From    string `json:"from"`
	To      string `json:"to"`
	Content string `json:"content"`
}
