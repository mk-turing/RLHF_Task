package main

import (
	"fmt"
	"log"
)

// Utility function to check for nil values and log an appropriate message
func checkNil(value interface{}, message string) error {
	if value == (*User)(nil) {
		return fmt.Errorf("%s is nil", message)
	}
	return nil
}

// Function that processes a user and returns a formatted string
func processUser(user *User) (string, error) {
	if err := checkNil(user, "user"); err != nil {
		return "", err
	}

	return fmt.Sprintf("User: %s, Age: %d", user.Name, user.Age), nil
}

// Function that processes a list of users and returns a formatted string
func processUsers(users []*User) (string, error) {
	if err := checkNil(users, "users"); err != nil {
		return "", err
	}

	result := "Users:\n"
	for _, user := range users {
		userString, err := processUser(user)
		if err != nil {
			return "", err
		}
		result += userString + "\n"
	}
	return result, nil
}

func main() {
	user1 := &User{"Alice", 25}
	user2 := (*User)(nil)
	users := []*User{user1, user2}

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

	resultUsers, errUsers := processUsers(users)
	if errUsers != nil {
		log.Printf("Error processing users: %v\n", errUsers)
	} else {
		fmt.Printf("%s", resultUsers)
	}
}

type User struct {
	Name string
	Age  int
}
