package main

import (
	"fmt"
	"github.com/pkg/errors"
)

type MyError struct {
	Message    string
	Service    string
	Method     string
	WrappedErr error
}

func (e *MyError) Error() string {
	return fmt.Sprintf("%s: %s in %s.%s", e.Message, e.WrappedErr.Error(), e.Service, e.Method)
}

func DoSomething() (result string, err error) {
	// Simulate an error scenario
	err = errors.Wrap(errors.New("internal error"), "failed to fetch data")
	return "", &MyError{
		Message:    "Something went wrong",
		Service:    "MyService",
		Method:     "DoSomething",
		WrappedErr: err,
	}
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
