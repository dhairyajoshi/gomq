package messages

import "time"

type Message struct {
	Id          string `json:"id"`
	EnqueuedAt  string `json:"enqueued_at"`
	CreatedAt   string `json:"created_at"`
	DeliveredAt time.Time
	Data        string `json:"data"`
}
