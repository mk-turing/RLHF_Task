package main

import (
	"fmt"
	"sync"
)

// User type represents a user in the collaboration system
type User string

// Document type represents the shared document
type Document []string

// Change represents a single change made to the document
type Change struct {
	User      User
	Operation string // "insert" or "delete"
	Position  int
	Text      string // For insert, the text to insert
}

// DocumentHandler is a higher-order function that takes a callback
// function as an argument. It is responsible for handling changes to the document.
type DocumentHandler func(func(Change))

// SynchronizeChanges creates a new DocumentHandler that synchronizes changes
// made to the document among multiple users.
func SynchronizeChanges() DocumentHandler {
	var document Document
	var changes []Change
	var mu sync.Mutex

	return func(onChange func(Change)) {
		for {
			// Wait for a change to be added
			change := <-changes

			mu.Lock()
			defer mu.Unlock()

			// Apply the change to the document
			applyChange(&document, change)

			// Broadcast the change to other users
			onChange(change)
		}
	}
}

// applyChange applies a single change to the document
func applyChange(document *Document, change Change) {
	switch change.Operation {
	case "insert":
		*document = append((*document)[:change.Position], change.Text)
		*document = append(*document, (*document)[change.Position:]...)
	case "delete":
		*document = append((*document)[:change.Position], (*document)[change.Position+1:]...)
	}
}

func main() {
	// Create multiple users
	users := []User{"Alice", "Bob", "Charlie"}

	// Create a new DocumentHandler to synchronize changes
	handler := SynchronizeChanges()

	// Start the synchronization process
	go handler(func(change Change) {
		fmt.Printf("User %s made change: %v\n", change.User, change)
	})

	// Simulate user actions
	for _, user := range users {
		go func(user User) {
			for i := 0; i < 5; i++ {
				position := i
				text := fmt.Sprintf("Hello from %s %d!", user, i)

				change := Change{
					User:      user,
					Operation: "insert",
					Position:  position,
					Text:      text,
				}

				// Send the change to the DocumentHandler for synchronization
				handler(func(ignore Change) {
					change
				})
			}
		}(user)
	}

	// Wait for all users to complete their actions
	// In a real application, we would use synchronization primitives to wait
	fmt.Scanln()
}
