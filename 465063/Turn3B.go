package main

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"os"
)

// LogError logs the error with additional context
func LogError(err error, context map[string]interface{}) {
	errorDetails := errors.Cause(err)
	// Format the error message with additional context
	errMsg := fmt.Sprintf("Error: %s", errorDetails.Error())
	for key, value := range context {
		errMsg = fmt.Sprintf("%s, %s: %v", errMsg, key, value)
	}
	log.Println(errMsg)
	// Optionally, you can use a logging library like logrus to add more features
	// logrus.WithFields(context).Error(errorDetails)
}

type MyCustomError struct {
	Message    string
	Service    string
	Method     string
	StackTrace string
}

func (e *MyCustomError) Error() string {
	return fmt.Sprintf("%s in %s.%s\nStack Trace: %s", e.Message, e.Service, e.Method, e.StackTrace)
}

// WrapErrorWithStackTrace wraps errors with stack trace and additional context
func WrapErrorWithStackTrace(err error, service string, method string) error {
	stackTrace := errors.WithStack(err).Error()
	return &MyCustomError{
		Message:    err.Error(),
		Service:    service,
		Method:     method,
		StackTrace: stackTrace,
	}
}

// Simulate a service that performs an operation that may fail
func DoSomething() (result string, err error) {
	// Simulate an error
	err = errors.New("internal server error")
	// Wrap the error with stack trace and additional context
	err = WrapErrorWithStackTrace(err, "MyService", "DoSomething")
	err = errors.Wrap(err, "while performing operation X")
	return "", err
}

func main() {
	// Set up the logger
	log.SetOutput(os.Stdout)

	// Call the service function
	result, err := DoSomething()
	if err != nil {
		// Log the error using the error logging service
		LogError(err, map[string]interface{}{
			"request_id": "1234",
			"user_id":    "user@example.com",
		})
		// Handle the error appropriately, e.g., send an error response
		fmt.Println("Error:", err.Error())
	} else {
		fmt.Println("Result:", result)
	}
}
