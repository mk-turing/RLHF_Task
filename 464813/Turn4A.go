package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	// Page size for paginated requests
	pageSize = 100
)

// Sample struct to represent a user
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Sample response struct that includes the user list and pagination details
type UserResponse struct {
	Users    []User `json:"users"`
	Total    int    `json:"total"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

func fetchUsers(page int) (*UserResponse, error) {
	// Rate limit the API requests
	rateLimiter <- struct{}{}
	defer func() { <-rateLimiter }()

	// Build the URL for the requested page
	url := fmt.Sprintf("https://api.example.com/users?page=%d&page_size=%d", page, pageSize)

	// Make a GET request to fetch users from an API
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		// Handle non-OK status code
		return nil, fmt.Errorf("received status %d", res.StatusCode)
	}

	// Decode the JSON response into a UserResponse struct
	var response UserResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
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

func fetchAndAddToQueue(page chan int, userPages chan []User, wg *sync.WaitGroup) {
	defer wg.Done()

	for pageNum := range page {
		response, err := fetchUsers(pageNum)
		if err != nil {
			log.Println("Error fetching page:", err)
			continue
		}
		userPages <- response.Users
	}
}

func processUsersConcurrently(users chan []User, wg *sync.WaitGroup) {
	defer wg.Done()

	for usersOnPage := range users {
		wg.Add(len(usersOnPage))

		for _, user := range usersOnPage {
			go func(user User) {
				processUser(user)
				wg.Done()
			}(user)
		}
	}
}

func processUser(user User) {
	fmt.Printf("Processing user %d: %s\n", user.ID, user.Name)
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
}

func handleUsersRequest(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup

	// Fetch users from the API with rate limiting and error handling
	userPages := make(chan []User)

	// Calculate the number of pages
	response, err := fetchUsers(1)
	if err != nil {
		log.Println("Error fetching users:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	numPages := response.Total / pageSize
	if response.Total%pageSize != 0 {
		numPages++
	}

	// Create a channel to fetch and add pages to the userPages channel
	page := make(chan int)

	// Create goroutines to fetch pages concurrently
	for i := 1; i <= numPages; i++ {
		page <- i
	}

	// Close the page channel once all page numbers are sent
	close(page)

	// Start goroutines to fetch and enqueue pages
	for i := 0; i < 5; i++ { // Adjust the number of fetching goroutines as needed
		wg.Add(1)
		go fetchAndAddToQueue(page, userPages, &wg)
	}

	// Start goroutines to process users concurrently from the queue
	for i := 0; i < 5; i++ { // Adjust the number of processing goroutines as needed
		wg.Add(1)
		go processUsersConcurrently(userPages, &wg)
	}

	// Wait for all fetching and processing goroutines to finish
	wg.Wait()

	// Write response to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Users processed concurrently with pagination!"))
}

func main() {
	http.HandleFunc("/process-users", handleUsersRequest)
	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
