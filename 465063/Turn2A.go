package main

import (
	"fmt"
	"github.com/pkg/errors"
)

// Define a custom error type
type MyCustomError struct {
	Message string
	Service string
	Method  string
}

func (e *MyCustomError) Error() string {
	return fmt.Sprintf("%s in %s.%s", e.Message, e.Service, e.Method)
}

// Function that does something
func DoSomething() (result string, err error) {
	// Simulate an error
	err = errors.Wrap(&MyCustomError{
		Message: "Something went wrong",
		Service: "MyService",
		Method:  "DoSomething",
	}, "while performing operation")
	return "", err
}

func main() {
	result, err := DoSomething()
	if err != nil {
		fmt.Println("Error:", err.Error())
		// Handle the error appropriately
	} else {
		fmt.Println("Result:", result)
	}
}
