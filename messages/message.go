package messages

type Message struct {
	Id         string `json:"id"`
	EnqueuedAt string `json:"enqueued_at"`
	CreatedAt  string `json:"created_at"`
	Data       []byte `json:"data"`
}
