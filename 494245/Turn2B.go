package main

import (
	"fmt"
	"sync"
	"time"
)

type Message struct {
	Content   string
	Timestamp time.Time
	Sender    string
}

type MessageStore struct {
	Messages []*Message
}

func NewMessageStore() *MessageStore {
	return &MessageStore{
		Messages: make([]*Message, 0),
	}
}

type User struct {
	Name  string
	Inbox chan *Message
}

var (
	users     sync.Map              // Stores user information
	broadcast = make(chan *Message) // Channel for broadcasting messages
)

func init() {
	// Start the message dispatcher
	go messageDispatcher()
}

func messageDispatcher() {
	for msg := range broadcast {
		// Broadcast message to all users
		users.Range(func(key, value interface{}) bool {
			user := value.(*User)
			user.Inbox <- msg
			return true
		})
	}
}

func addUser(name string) {
	user := &User{
		Name:  name,
		Inbox: make(chan *Message),
	}
	users.Store(name, user)
	go handleUserMessages(user)
}

func handleUserMessages(user *User) {
	for msg := range user.Inbox {
		fmt.Printf("%s: %s - %s\n", msg.Sender, msg.Timestamp.Format("2006-01-02 15:04:05"), msg.Content)
	}
}

func sendMessage(sender string, content string, recipient string) {
	msg := &Message{
		Content:   content,
		Timestamp: time.Now(),
		Sender:    sender,
	}

	if recipient == "" {
		// Broadcast message to all users
		broadcast <- msg
	} else {
		// Send private message to the recipient
		if recvUser, ok := users.Load(recipient); ok {
			recvUser.(*User).Inbox <- msg
		} else {
			fmt.Printf("User %s not found.\n", recipient)
		}
	}
}

func displayUsers() {
	fmt.Println("Users:")
	users.Range(func(key, value interface{}) bool {
		fmt.Println(key)
		return true
	})
}

func main() {
	fmt.Println("Welcome to the Enhanced Messaging Application!")

	for {
		fmt.Print("Command: ")
		var cmd string
		fmt.Scanln(&cmd)

		switch cmd {
		case "add":
			fmt.Print("Enter user name: ")
			var userName string
			fmt.Scanln(&userName)
			addUser(userName)
			fmt.Println("User added.")
		case "send":
			fmt.Print("Enter sender name: ")
			var sender string
			fmt.Scanln(&sender)
			fmt.Print("Enter message content: ")
			var messageContent string
			fmt.Scanln(&messageContent)
			fmt.Print("Enter recipient (empty for broadcast): ")
			var recipient string
			fmt.Scanln(&recipient)
			sendMessage(sender, messageContent, recipient)
			fmt.Println("Message sent.")
		case "display":
			displayUsers()
		case "exit":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid command. Try add, send, display, or exit.")
		}
	}
}
