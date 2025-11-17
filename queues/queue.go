package queues

import (
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/dhairyajoshi/gomq/io"
	"github.com/dhairyajoshi/gomq/messages"
	"github.com/dhairyajoshi/gomq/parsers"
)

type Queue interface {
	Enqueue(message messages.Message) bool
	Consume() (messages.Message, bool)
	getDelivered() []messages.Message
	getName() string
	requeueMessage(idx int) bool
	Subscribe(*io.IOHandler) bool
}

type DurableQueue struct {
	name        string
	messages    []messages.Message
	delivered   []messages.Message
	subscribers []*io.IOHandler
	lock        sync.Mutex
}

func (q *DurableQueue) getName() string {
	return q.name
}

func (q *DurableQueue) Enqueue(message messages.Message) bool {
	q.lock.Lock()
	defer q.lock.Unlock()
	if len(q.subscribers) > 0 {
		closedConns := []*io.IOHandler{}
		for _, sub := range q.subscribers {
			_, err := (*sub).Write(parsers.ServerResponse{Type: "message", Data: message, SendNext: false, Close: false})
			if err != nil {
				fmt.Println("error sending message to subscriber: ", err.Error())
				(*sub).Write(parsers.ServerResponse{Type: "server_response", Data: "Closing connection!", SendNext: false, Close: true})
				(*sub).Close()
				closedConns = append(closedConns, sub)
			}
		}
		for _, closedConn := range closedConns {
			q.subscribers = slices.DeleteFunc(q.subscribers, func(sub *io.IOHandler) bool { return sub == closedConn })
			fmt.Println("removed closed connection from subscribers!")
		}
	} else {
		q.messages = append(q.messages, message)
	}
	return true
}

func (q *DurableQueue) Consume() (messages.Message, bool) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if len(q.messages) == 0 {
		return messages.Message{}, false
	}
	message := q.messages[0]
	q.messages = q.messages[1:]
	message.DeliveredAt = time.Now()
	q.delivered = append(q.delivered, message)
	return message, true
}

func (q *DurableQueue) getDelivered() []messages.Message {
	q.lock.Lock()
	defer q.lock.Unlock()
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

func (q *DurableQueue) Subscribe(conn *io.IOHandler) bool {
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

var (
	QueueStore = map[string]Queue{}
	QueueLock  sync.RWMutex
)

func NewDurableQueue(name string) (Queue, bool) {
	QueueLock.Lock()
	defer QueueLock.Unlock()
	queue := &DurableQueue{name: name}
	QueueStore[name] = queue
	return queue, true
}

func GetQueue(name string) (*Queue, error) {
	QueueLock.RLock()
	defer QueueLock.RUnlock()
	queue, found := QueueStore[name]
	if !found {
		return nil, NoSuchQueueError{}
	}
	return &queue, nil
}

func GetOrCreateDurableQueue(name string) (Queue, bool) {
	QueueLock.RLock()
	queue, found := QueueStore[name]
	if !found {
		QueueLock.RUnlock()
		return NewDurableQueue(name)
	}
	return queue, false
}

func MonitorQueues() {
	for {
		QueueLock.RLock()
		for k := range QueueStore {
			queue := QueueStore[k]
			for idx, message := range queue.getDelivered() {
				if time.Since(message.DeliveredAt) >= 10*time.Second {
					fmt.Println("requeueing un-acked message ", message.Id, " in queue", queue.getName())
					queue.requeueMessage(idx)
				}
			}
		}
		QueueLock.RUnlock()
	}
}
