package main

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

// Sample struct to represent a user
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Define a semaphore
type Semaphore struct {
	count int32
	mu    sync.Mutex
}

func NewSemaphore(count int) *Semaphore {
	return &Semaphore{count: int32(count)}
}

func (s *Semaphore) Acquire() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for s.count <= 0 {
		time.Sleep(time.Millisecond * 100) // Wait briefly if the semaphore is zero
	}

	atomic.AddInt32(&s.count, -1)
}

func (s *Semaphore) Release() {
	s.mu.Lock()
	defer s.mu.Unlock()

	atomic.AddInt32(&s.count, 1)
}

func fetchUsers(sem *Semaphore) ([]User, error) {
	// Simulate making a GET request to fetch users from an API
	// This function should include actual HTTP request logic
	// For this example, we'll simulate a successful response

	time.Sleep(time.Duration(100) * time.Millisecond) // Simulate request latency

	users := []User{
		{ID: 1, Name: "Alice", Email: "alice@example.com"},
		{ID: 2, Name: "Bob", Email: "bob@example.com"},
		{ID: 3, Name: "Charlie", Email: "charlie@example.com"},
		{ID: 4, Name: "Mayank", Email: "mayank@example.com"},
	}

	return users, nil
}

func processUser(user User) {
	fmt.Printf("Processing user %d: %s\n", user.ID, user.Name)
}

func processUsersConcurrently(sem *Semaphore, users []User) {
	// Create a wait group to synchronize goroutines
	var wg sync.WaitGroup

	// Specify a fixed number of goroutines to process the users concurrently
	numGoroutines := 2 // For example, we use 2 goroutines
	chunkSize := len(users) / numGoroutines

	// Split users into chunks and send each chunk to a goroutine
	for i := 0; i < numGoroutines; i++ {
		start := i * chunkSize
		end := min(start+chunkSize, len(users))

		wg.Add(1)
		go func(users []User) {
			defer wg.Done()
			for _, user := range users {
				sem.Acquire() // Acquire semaphore before processing user
				processUser(user)
				sem.Release() // Release semaphore after processing user
			}
		}(users[start:end])
	}

	// Wait for all goroutines to finish
	wg.Wait()
	fmt.Println("All users processed")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	// Initialize a semaphore with a limit of 4 concurrent requests
	sem := NewSemaphore(4)

	// Simulate fetching users from the API
	users, err := fetchUsers(sem)
	if err != nil {
		log.Println(err)
		return
	}

	// Process users in parallel with rate-limiting
	processUsersConcurrently(sem, users)

	fmt.Println("Users processed concurrently with rate-limiting")
}
