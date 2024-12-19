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
	apiRequestLimit = 10
	maxBurstSize    = 5
	baseURL         = "https://api.example.com/users?page=%d"
)

// Sample struct to represent a user
type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
}

type APIResponse struct {
	Data  []User `json:"data"`
	NextPage int `json:"next_page"`
}

func fetchUsers(page int) (*APIResponse, error) {
	// Rate limit the API requests
	rateLimiter <- struct{}{}
	defer func() { <-rateLimiter }()

	url := fmt.Sprintf(baseURL, page)
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
		return nil, fmt.Errorf("received status %d", res.StatusCode)
	}

	var response APIResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Create a channel with a buffer to handle rate limiting
var rateLimiter = make(chan struct{}, maxBurstSize)

func init() {
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
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
}

func processPage(page int, pageProcessed chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	response, err := fetchUsers(page)
	if err != nil {
		log.Println("Error fetching page", page, ":", err)
		return
	}
	for _, user := range response.Data {
		go processUser(user, &sync.WaitGroup{})
	}
	pageProcessed <- page
}

func processUsersConcurrently(totalPages int) {
	var wg sync.WaitGroup
	pageProcessed := make(chan int)

	// Start a goroutine to process each page
	for page := 1; page <= totalPages; page++ {
		wg.Add(1)
		go processPage(page, pageProcessed, &wg)
	}

	// Wait for all pages to be processed
	go func() {
		wg.Wait()
		close(pageProcessed)
	}()

	// Display processed pages as they complete
	for page := range pageProcessed {
		fmt.Println("Page", page, "processed")
	}
}

func handleUsersRequest(w http.ResponseWriter, r *http.Request) {
	// Assume the total number of pages is known or can be fetched from the API
	totalPages := 5
	processUsersConcurrently(totalPages)

	w.WriteHeader(http.StatusOK)