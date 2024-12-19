package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

func handleUsersRequest(w http.ResponseWriter, r *http.Request) {
	// Fetch users from the API
	users, err := fetchUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Process users with range
	for _, user := range users {
		fmt.Println(fmt.Sprintf("User ID: %d, Name: %s, Email: %s", user.ID, user.Name, user.Email))
	}

	// Write response to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Users processed successfully!"))
}

func main() {
	http.HandleFunc("/process-users", handleUsersRequest)
	fmt.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}
