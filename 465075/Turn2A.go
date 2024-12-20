package main

import (
	"fmt"
	"log"
)

// Custom error type for handling nil values
type NilError struct {
	Message string
}

func (e *NilError) Error() string {
	return e.Message
}

// Function that processes a user and returns a formatted string
func processUser(user *User) (string, error) {
	if user == nil {
		return "", &NilError{Message: "user cannot be nil"}
	}

	return fmt.Sprintf("User: %s, Age: %d", user.Name, user.Age), nil
}

func main() {
	user1 := &User{"Alice", 25}
	user2 := (*User)(nil)

	result1, err1 := processUser(user1)
	if err1 != nil {
		log.Printf("Error processing user1: %v\n", err1)
	} else {
		fmt.Printf("User1: %s\n", result1)
	}

	result2, err2 := processUser(user2)
	if err2 != nil {
		log.Printf("Error processing user2: %v\n", err2)
	} else {
		fmt.Printf("User2: %s\n", result2)
	}
}

type User struct {
	Name string
	Age  int
}
