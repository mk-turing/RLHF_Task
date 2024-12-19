package main

import (
	"fmt"
)

type MyError struct {
	Message string
}

func (e *MyError) Error() string {
	return e.Message
}

func DoSomething() (result string, err error) {
	// Simulate an error scenario
	return "", &MyError{Message: "Something went wrong"}
}

func main() {
	result, err := DoSomething()
	if err != nil {
		fmt.Println("Error:", err.Error())
		// Handle the error appropriately, e.g., log it, send an error response, or retry
	} else {
		fmt.Println("Result:", result)
	}
}
