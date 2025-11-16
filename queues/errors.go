package queues

type NoSuchQueueError struct{}

func (NoSuchQueueError) Error() string {
	return "No such queue exists!"
}
