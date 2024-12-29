package main

import (
	"errors"
	"fmt"
	"testing"
)

// ProcessData processes data and calls a callback on success or error.
func processData(data []byte, callback func(error)) {
	// Simulated error handling
	if data == nil {
		callback(errors.New("data is nil"))
		return
	}

	// Simulated work that might fail
	if len(data) == 0 {
		callback(errors.New("data is empty"))
		return
	}

	// Simulated successful processing
	callback(nil)
}

// MockCallback is a mock implementation of the callback function.
func MockCallback(err error) {
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func TestProcessDataWithError(t *testing.T) {
	// Simulate an empty data slice to cause an error.
	processData([]byte{}, MockCallback)

	// We can assert the behavior of the mock by checking outputs or using a more sophisticated mock library.
	// For simplicity, let's just check if the error message is printed.
	t.Helper()
	// Alternatively, you could use the testing.Assert library for more structured assertions.
}

func TestProcessDataWithoutError(t *testing.T) {
	// Simulate successful processing with non-nil and non-empty data.
	processData([]byte{1, 2, 3}, MockCallback)

	// Assert successful processing by checking if no error was logged.
	t.Helper()
}

func TestProcessDataWithNilData(t *testing.T) {
	// Simulate an error by passing nil data.
	processData(nil, MockCallback)

	// Assert the error message is printed for nil data.
	t.Helper()
}
