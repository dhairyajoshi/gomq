package commands

import (
	"fmt"
	"time"

	"github.com/dhairyajoshi/gomq/io"
	"github.com/dhairyajoshi/gomq/messages"
	"github.com/dhairyajoshi/gomq/parsers"
	"github.com/dhairyajoshi/gomq/queues"
	"github.com/google/uuid"
)

func createQueue(args ...any) parsers.ServerResponse {
	name, ok := args[1].(string)
	if !ok {
		fmt.Println("Name not a valid string: ", name)
		return parsers.ServerResponse{Data: "Name not a valid string", SendNext: true, Close: false, Type: "server_response"}
	}
	_, created := queues.GetOrCreateDurableQueue(name)
	if created {
		return parsers.ServerResponse{Data: "Queue created Successfully", SendNext: true, Close: false, Type: "server_response"}
	}
	return parsers.ServerResponse{Data: "Queue already exists", SendNext: true, Close: false, Type: "server_response"}
}

func publishMessage(args ...any) parsers.ServerResponse {
	queueName, ok := args[1].(string)
	if !ok {
		fmt.Println("invalid queue name ", queueName)
		return parsers.ServerResponse{Data: "invalid queue name", SendNext: true, Close: false, Type: "server_response"}
	}
	message, ok := args[2].(string)
	if !ok {
		fmt.Println("Invalid message ", message)
		return parsers.ServerResponse{Data: "invalid message", SendNext: true, Close: false, Type: "server_response"}
	}
	queue, _ := queues.GetOrCreateDurableQueue(queueName)
	ok = queue.Enqueue(messages.Message{Id: uuid.New().String(), Data: message, EnqueuedAt: time.Now().String()})
	if !ok {
		return parsers.ServerResponse{Data: "Couldn't enqueue message", SendNext: true, Close: false, Type: "server_response"}
	}
	return parsers.ServerResponse{Data: "Message enqueued Successfully", SendNext: true, Close: false, Type: "server_response"}
}

func consumeMessage(args ...any) parsers.ServerResponse {
	queueName, ok := args[1].(string)
	if !ok {
		fmt.Println("invalid queue name ", queueName)
		return parsers.ServerResponse{Data: "invalid queue name", SendNext: true, Close: false, Type: "server_response"}
	}
	queue, _ := queues.GetOrCreateDurableQueue(queueName)
	message := queue.Consume()
	if !ok {
		return parsers.ServerResponse{Data: "Couldn't consume message", SendNext: true, Close: false, Type: "server_response"}
	}
	return parsers.ServerResponse{Data: message, SendNext: true, Close: false, Type: "message"}
}

func subscribeQueue(args ...any) parsers.ServerResponse {
	conn, ok := args[0].(*io.IOHandler)
	if !ok {
		fmt.Println("Didn't receive valid connection!")
		return parsers.ServerResponse{Data: "Didn't receive valid connection!", SendNext: true, Close: false, Type: "server_response"}
	}
	queueName, ok := args[1].(string)
	if !ok {
		fmt.Println("invalid queue name ", queueName)
		return parsers.ServerResponse{Data: "invalid queue name", SendNext: true, Close: false, Type: "server_response"}
	}
	queue, _ := queues.GetOrCreateDurableQueue(queueName)
	ok = queue.Subscribe(conn)
	if !ok {
		return parsers.ServerResponse{Data: "Couldn't subscribe to queue", SendNext: true, Close: false, Type: "server_response"}
	}
	return parsers.ServerResponse{Data: "Subscribed to queue, you will receive messages as they're published!", SendNext: false, Close: false, Type: "server_response"}
}

var queueCommands = map[string]func(args ...any) parsers.ServerResponse{
	"create-queue":    createQueue,
	"publish-message": publishMessage,
	"consume-message": consumeMessage,
	"subscribe-queue": subscribeQueue,
}
