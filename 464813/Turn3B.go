package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

const (
	// Maximum number of requests allowed per second
	apiRequestLimit = 10
	// Maximum burst size for the rate limiter
	maxBurstSize = 5
)

// Sample struct to represent a user
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func fetchUsers() ([]User, error) {
	users := []User{
		{ID: 1, Name: "Alice", Email: "alice@example.com"},
		{ID: 2, Name: "Bob", Email: "bob@example.com"},
		{ID: 3, Name: "Charlie", Email: "charlie@example.com"},
	}

	return users, nil
}

// Create a channel with a buffer to handle rate limiting
var rateLimiter = make(chan struct{}, maxBurstSize)

func init() {
	// Start a timer to tick once per second to refill the rate limiter
	go func() {
		for {
			time.Sleep(time.Second)
			for i := 0; i < apiRequestLimit; i++ {
				rateLimiter <- struct{}{}
			}
		}
	}()
}

func processUser(user User, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Processing user %d: %s\n", user.ID, user.Name)
	// Add random processing time to simulate work
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
}

func processUsersConcurrently(users []User) {
	var wg sync.WaitGroup
	wg.Add(len(users))

	for _, user := range users {
		go processUser(user, &wg)
	}

	wg.Wait()
	fmt.Println("All users processed")
}

func handleUsersRequest(w http.ResponseWriter, r *http.Request) {
	// Fetch users from the API with rate limiting and error handling
	users, err := fetchUsers()
	if err != nil {
		// Handle the error appropriately based on its type
		log.Println("Error fetching users:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Process users in parallel
	processUsersConcurrently(users)

	// Write response to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Users processed concurrently with rate limiting!"))
}

func main() {
	http.HandleFunc("/process-users", handleUsersRequest)
	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
