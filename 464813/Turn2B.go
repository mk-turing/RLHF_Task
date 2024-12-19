package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

// Sample struct to represent a user
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func fetchUsers(url string) ([]User, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var users []User
	err = json.Unmarshal(body, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func processUser(user User, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println(fmt.Sprintf("Processing User ID: %d, Name: %s, Email: %s", user.ID, user.Name, user.Email))
}

func handleUsersRequest(w http.ResponseWriter, r *http.Request) {
	urls := []string{
		"https://api.example.com/users1",
		"https://api.example.com/users2",
		"https://api.example.com/users3", // Add more URLs as needed
	}

	var wg sync.WaitGroup
	wg.Add(len(urls))

	for _, url := range urls {
		go func(url string) {
			defer wg.Done()
			users, err := fetchUsers(url)
			if err != nil {
				// Handle error, e.g., log it or send to error channel
				return
			}
			for _, user := range users {
				go processUser(user, &wg)
			}
		}(url)
	}

	wg.Wait()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("All users processed successfully!"))
}

func main() {
	http.HandleFunc("/process-users", handleUsersRequest)
	fmt.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}
