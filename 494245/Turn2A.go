package main

import (
	"fmt"
	"sync"
	"time"
)

// User struct holds user information and their message channel
type User struct {
	Name   string
	MsgCh  chan *Message
	active bool
}

// Message struct contains the message content and timestamp
type Message struct {
	Content   string
	Timestamp time.Time
	Sender    string
}

// MessageStore contains a slice of users and some utility functions
type MessageStore struct {
	Users         map[string]*User
	broadcastCh   chan *Message
	userLock      sync.RWMutex
	broadcastLock sync.Mutex
}

// NewMessageStore returns a new instance of MessageStore
func NewMessageStore() *MessageStore {
	return &MessageStore{
		Users:         make(map[string]*User),
		broadcastCh:   make(chan *Message),
		userLock:      sync.RWMutex{},
		broadcastLock: sync.Mutex{},
	}
}

// AddUser adds a new user to the store
func (ms *MessageStore) AddUser(name string) {
	ms.userLock.Lock()
	defer ms.userLock.Unlock()
	if _, exists := ms.Users[name]; exists {
		fmt.Println("User already exists.")
		return
	}
	ms.Users[name] = &User{
		Name:   name,
		MsgCh:  make(chan *Message),
		active: true,
	}
	go ms.userListener(name)
	fmt.Printf("User %s added.\n", name)
}

// RemoveUser removes a user from the store
func (ms *MessageStore) RemoveUser(name string) {
	ms.userLock.Lock()
	defer ms.userLock.Unlock()
	if user, exists := ms.Users[name]; exists {
		user.active = false
		close(user.MsgCh)
		delete(ms.Users, name)
		fmt.Printf("User %s removed.\n", name)
	} else {
		fmt.Println("User not found.")
	}
}

// SendMessage sends a message to a specific user
func (ms *MessageStore) SendMessage(sender, receiver, content string) {
	ms.userLock.RLock()
	defer ms.userLock.RUnlock()
	if receiverUser, exists := ms.Users[receiver]; exists && receiverUser.active {
		receiverUser.MsgCh <- &Message{
			Content:   content,
			Timestamp: time.Now(),
			Sender:    sender,
		}
		fmt.Printf("Message sent from %s to %s.\n", sender, receiver)
	} else {
		fmt.Printf("User %s not found or not active.\n", receiver)
	}
}

// BroadcastMessage sends a message to all active users
func (ms *MessageStore) BroadcastMessage(sender, content string) {
	msg := &Message{
		Content:   content,
		Timestamp: time.Now(),
		Sender:    sender,
	}

	ms.broadcastLock.Lock()
	ms.userLock.RLock()
	defer ms.userLock.RUnlock()
	defer ms.broadcastLock.Unlock()

	for _, user := range ms.Users {
		if user.active {
			user.MsgCh <- msg
		}
	}

	ms.broadcastCh <- msg
}

// userListener listens for messages on a user's channel and prints them
func (ms *MessageStore) userListener(name string) {
	for msg := range ms.Users[name].MsgCh {
		fmt.Printf("%s (%s): %s\n", msg.Timestamp.Format("2006-01-02 15:04:05"), msg.Sender, msg.Content)
	}
}

func main() {
	store := NewMessageStore()
	fmt.Println("Welcome to the Extended Messaging Application!")

	go store.broadcastListener()

	for {
		fmt.Print("Command: ")
		var cmd string
		fmt.Scanln(&cmd)
		switch cmd {
		case "add":
			fmt.Print("Enter user name: ")
			var name string
			fmt.Scanln(&name)
			store.AddUser(name)
		case "remove":
			fmt.Print("Enter user name to remove: ")
			var name string
			fmt.Scanln(&name)
			store.RemoveUser(name)
		case "send":
			fmt.Print("Enter sender name: ")
			var sender string
			fmt.Scanln(&sender)
			fmt.Print("Enter receiver name: ")
			var receiver string
			fmt.Scanln(&receiver)
			fmt.Print("Enter message: ")
			var content string
			fmt.Scanln(&content)
			store.SendMessage(sender, receiver, content)
		case "broadcast":
			fmt.Print("Enter sender name: ")
			var sender string
			fmt.Scanln(&sender)
			fmt.Print("Enter message: ")
			var content string
			fmt.Scanln(&content)
			store.BroadcastMessage(sender, content)
		case "exit":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid command. Try add, remove, send, broadcast, or exit.")
		}
	}
}

// broadcastListener listens for broadcast messages and prints them
func (ms *MessageStore) broadcastListener() {
	for msg := range ms.broadcastCh {
		fmt.Printf("%s (%s): %s\n", msg.Timestamp.Format("2006-01-02 15:04:05"), msg.Sender, msg.Content)
	}
}
