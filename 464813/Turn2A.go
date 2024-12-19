package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"sync"
)

// Sample struct to represent a user
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func fetchUsers() ([]User, error) {
	// Make a GET request to fetch users from an API
	res, err := http.Get("https://api.example.com/users")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Decode the JSON response into a slice of User structs
	var users []User
	err = json.Unmarshal(body, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func processUser(user User, wg *sync.WaitGroup) {
	defer wg.Done()
	// Simulate processing a user
	fmt.Printf("Processing user %d: %s\n", user.ID, user.Name)
}

func processUsersConcurrently(users []User) {
	// Create a wait group to synchronize goroutines
	var wg sync.WaitGroup
	// Calculate the number of goroutines to use
	numGoroutines := int(math.Ceil(float64(len(users)) / float64(8))) // Adjust the number as needed
	chunkSize := len(users) / numGoroutines

	// Split users into chunks and send each chunk to a goroutine
	for i := 0; i < numGoroutines; i++ {
		start := i * chunkSize
		end := min(start+chunkSize, len(users))

		wg.Add(end - start)
		go func(users []User) {
			for _, user := range users {
				processUser(user, &wg)
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

func handleUsersRequest(w http.ResponseWriter, r *http.Request) {
	// Fetch users from the API
	users, err := fetchUsers()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Process users in parallel
	processUsersConcurrently(users)

	// Write response to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Users processed concurrently!"))
}

func main() {
	http.HandleFunc("/process-users", handleUsersRequest)
	fmt.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println(err)
	}
}
