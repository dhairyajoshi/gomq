package queues

import (
	"fmt"
	"net"
	"slices"
	"sync"
	"time"

	"github.com/dhairyajoshi/gomq/messages"
)

type Queue interface {
	Enqueue(message messages.Message) bool
	Consume() messages.Message
	getDelivered() []messages.Message
	getName() string
	requeueMessage(idx int) bool
	Subscribe(*net.Conn) bool
}

type DurableQueue struct {
	name        string
	messages    []messages.Message
	delivered   []messages.Message
	subscribers []*net.Conn
	lock        sync.Mutex
}

func (q *DurableQueue) getName() string {
	return q.name
}

func (q *DurableQueue) Enqueue(message messages.Message) bool {
	if len(q.subscribers) > 0 {
		q.lock.Lock()
		for _, sub := range q.subscribers {
			_, err := (*sub).Write([]byte(message.Data))
			if err != nil {
				fmt.Println("error sending message to subscriber: ", err.Error())
				(*sub).Close()
			}
		}
	} else {
		q.messages = append(q.messages, message)
	}
	return true
}

func (q *DurableQueue) Consume() messages.Message {
	q.lock.Lock()
	message := q.messages[0]
	q.messages = q.messages[1:]
	message.DeliveredAt = time.Now()
	q.delivered = append(q.delivered, message)
	q.lock.Unlock()
	return message
}

func (q *DurableQueue) getDelivered() []messages.Message {
	return q.delivered
}

func (q *DurableQueue) requeueMessage(idx int) bool {
	q.lock.Lock()
	message := q.delivered[idx]
	q.messages = append(q.messages, message)
	q.delivered = slices.DeleteFunc(q.delivered, func(m messages.Message) bool { return m.Id == message.Id })
	q.lock.Unlock()
	return true
}

func (q *DurableQueue) Subscribe(conn *net.Conn) bool {
	q.lock.Lock()
	q.subscribers = append(q.subscribers, conn)
	q.lock.Unlock()
	return true
}

func (q *DurableQueue) Ack(message_id string) {
	q.lock.Lock()
	q.delivered = slices.DeleteFunc(q.delivered, func(e messages.Message) bool {
		return e.Id == message_id
	})
	q.lock.Unlock()
}

var QueueStore = map[string]Queue{}

func NewDurableQueue(name string) (Queue, bool) {
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

func GetOrCreateDurableQueue(name string) (Queue, bool) {
	queue, found := QueueStore[name]
	if !found {
		return NewDurableQueue(name)
	}
	return queue, false
}

func MonitorQueues() {
	for {
		for k := range QueueStore {
			queue := QueueStore[k]
			for idx, message := range queue.getDelivered() {
				if time.Since(message.DeliveredAt) >= 10*time.Second {
					fmt.Println("requeueing un-acked message ", message.Id, " in queue", queue.getName())
					queue.requeueMessage(idx)
				}
			}
		}
	}
}
