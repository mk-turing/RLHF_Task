package main

import (
	"fmt"
	"strings"
	"time"
)

// Message struct contains the message content and timestamp
type Message struct {
	Content   string
	Timestamp time.Time
}

// MessageStore contains a slice of messages and some utility functions
type MessageStore struct {
	Messages []*Message
}

// NewMessageStore returns a new instance of MessageStore
func NewMessageStore() *MessageStore {
	return &MessageStore{
		Messages: make([]*Message, 0),
	}
}

// AddMessage adds a new message to the store
func (ms *MessageStore) AddMessage(content string) {
	newMessage := &Message{
		Content:   strings.TrimSpace(content),
		Timestamp: time.Now(),
	}
	ms.Messages = append(ms.Messages, newMessage)
}

// DisplayLastNMessages displays the last N messages
func (ms *MessageStore) DisplayLastNMessages(n int) {
	if n <= 0 {
		fmt.Println("Invalid number of messages to display.")
		return
	}

	lastNMessages := ms.Messages[len(ms.Messages)-n:]
	for _, message := range lastNMessages {
		fmt.Printf("%s - %s\n", message.Timestamp.Format("2006-01-02 15:04:05"), message.Content)
	}
}

// ClearMessages clears all messages from the store
func (ms *MessageStore) ClearMessages() {
	ms.Messages = ms.Messages[:0]
}

func main() {
	store := NewMessageStore()
	fmt.Println("Welcome to the Basic Messaging Application!")

	for {
		fmt.Print("Command: ")
		var cmd string
		fmt.Scanln(&cmd)
		switch cmd {
		case "add":
			fmt.Print("Enter message: ")
			var messageContent string
			fmt.Scanln(&messageContent)
			store.AddMessage(messageContent)
			fmt.Println("Message added.")
		case "display":
			fmt.Print("Enter number of messages to display: ")
			var n int
			fmt.Scanln(&n)
			store.DisplayLastNMessages(n)
		case "clear":
			store.ClearMessages()
			fmt.Println("Messages cleared.")
		case "exit":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid command. Try add, display, clear, or exit.")
		}
	}
}
