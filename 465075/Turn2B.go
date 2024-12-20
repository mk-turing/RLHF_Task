package main

import (
	"fmt"
)

type User struct {
	Name string
	Age  int
}

// Define a custom error type
type NilUserError struct {
	msg string
}

func (e *NilUserError) Error() string {
	return fmt.Sprintf("NilUserError: %s", e.msg)
}
func processUser(user *User) error {
	if user == nil {
		return &NilUserError{"user parameter should not be nil"}
	}
	// Perform your processing here
	return nil
}
func main() {
	user1 := &User{"Alice", 25}
	user2 := (*User)(nil) // nil value
	err := processUser(user1)
	if err != nil {
		// Handle the error
		fmt.Println("Error:", err)
		return
	}
	err = processUser(user2)
	if err != nil {
		// Handle the error
		fmt.Println("Error:", err)
		return
	}
}
