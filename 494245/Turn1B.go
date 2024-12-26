package main

import (
	"fmt"
	"sync"
	"time"
)

type Message struct {
	Sender string
	Text   string
}

type MessageQueue struct {
	messages []*Message
	mu       sync.Mutex
}

func (q *MessageQueue) Push(message *Message) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.messages = append(q.messages, message)
}

func (q *MessageQueue) Pop() *Message {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.messages) == 0 {
		return nil
	}

	message := q.messages[0]
	q.messages = q.messages[1:]
	return message
}

func (q *MessageQueue) IsEmpty() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.messages) == 0
}

func main() {
	messageQueue := NewMessageQueue()

	go func() {
		for {
			message := messageQueue.Pop()
			if message != nil {
				fmt.Printf("Received: %s: %s\n", message.Sender, message.Text)
			}
		}
	}()

	sender := "Alice"
	message := "Hello, Bob!"
	messageQueue.Push(&Message{Sender: sender, Text: message})

	sender = "Bob"
	message = "Hi, Alice!"
	messageQueue.Push(&Message{Sender: sender, Text: message})

	// Allow the receiver goroutine to process the messages
	<-time.After(time.Second)
}
