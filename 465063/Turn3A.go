package main

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
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

// Mock error logging service
type Logger interface {
	LogError(err error, service string, method string)
}

type MockLogger struct{}

func (l *MockLogger) LogError(err error, service string, method string) {
	log.Printf("Service: %s, Method: %s, Error: %v", service, method, err)
}

// Function that does something
func DoSomething(logger Logger) (result string, err error) {
	// Simulate an error
	err = errors.Wrap(&MyCustomError{
		Message: "Something went wrong",
		Service: "MyService",
		Method:  "DoSomething",
	}, "while performing operation")

	// Log the error
	logger.LogError(err, "MyService", "DoSomething")

	return "", err
}

func main() {
	logger := &MockLogger{}

	result, err := DoSomething(logger)
	if err != nil {
		fmt.Println("Error:", err.Error())
		// Handle the error appropriately
	} else {
		fmt.Println("Result:", result)
	}
}
