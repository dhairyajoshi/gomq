package commands

import (
	"fmt"
	"time"

	"github.com/dhairyajoshi/gomq/messages"
	"github.com/dhairyajoshi/gomq/queues"
	"github.com/google/uuid"
)

func createQueue(args ...any) (bool, string) {
	name, ok := args[0].(string)
	if !ok {
		fmt.Println("Name not a valid string: ", name)
		return false, "Name not a valid string: "
	}
	_, created := queues.NewQueue(name)
	if created {
		return false, "Queue created Successfuly"
	}
	return false, "Couldn't create queue"
}

func publishMessage(args ...any) (bool, string) {
	queueName, ok := args[0].(string)
	if !ok {
		fmt.Println("invalid queue name ", queueName)
		return false, "invalid queue name"
	}
	message, ok := args[1].(string)
	if !ok {
		fmt.Println("Invalid message ", message)
		return false, "invalid message"
	}
	queue, _ := queues.GetOrCreateQueue(queueName)
	ok = queue.Enqueue(messages.Message{Id: uuid.New().String(), Data: []byte(message), EnqueuedAt: time.Now().String()})
	if !ok {
		return false, "Couldn't enqueue message"
	}
	return false, "Message enqueued Successfuly"
}

func consumeMessage(args ...any) (bool, string) {
	queueName, ok := args[0].(string)
	if !ok {
		fmt.Println("invalid queue name ", queueName)
		return false, "invalid queue name"
	}
	queue, _ := queues.GetOrCreateQueue(queueName)
	message := queue.Consume()
	if !ok {
		return false, "Couldn't consume message"
	}
	return false, string(message.Data)
}

var queueCommands = map[string]func(args ...any) (bool, string){
	"create-queue":    createQueue,
	"publish-message": publishMessage,
	"consume-message": consumeMessage,
}
