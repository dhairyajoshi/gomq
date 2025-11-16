package queues

import (
	"slices"
	"sync"

	"github.com/dhairyajoshi/gomq/messages"
)

type Queue interface {
	Enqueue(message messages.Message) bool

	Consume() messages.Message
}

type DurableQueue struct {
	name      string
	messages  []messages.Message
	delivered []messages.Message
	lock      sync.Mutex
}

func (q *DurableQueue) Enqueue(message messages.Message) bool {
	q.messages = append(q.messages, message)
	return true
}

func (q *DurableQueue) Consume() messages.Message {
	q.lock.Lock()
	message := q.messages[0]
	q.messages = q.messages[1:]
	q.delivered = append(q.delivered, message)
	q.lock.Unlock()
	return message
}

func (q *DurableQueue) Ack(message_id string) {
	q.lock.Lock()
	q.delivered = slices.DeleteFunc(q.delivered, func(e messages.Message) bool {
		return e.Id == message_id
	})
	q.lock.Unlock()
}

var QueueStore = map[string]Queue{}

func NewQueue(name string) (Queue, bool) {
	queue := &DurableQueue{name: name}
	QueueStore[name] = queue
	return queue, true
}

func GetQueue(name string) (*Queue, error) {
	queue, found := QueueStore[name]
	if !found {
		return nil, NoSuchQueueError{}
	}
	return &queue, nil
}

func GetOrCreateQueue(name string) (Queue, bool) {
	queue, found := QueueStore[name]
	if !found {
		return NewQueue(name)
	}
	return queue, false
}
