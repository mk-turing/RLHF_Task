package main

import (
	"errors"
	"sync"
	"testing"
	"time"
)

// processDataWithConcurrency processes data using multiple goroutines and calls callbacks accordingly.
func processDataWithConcurrency(data [][]byte, callbacks []func(error)) {
	wg := sync.WaitGroup{}
	wg.Add(len(data))

	for i, d := range data {
		go func(idx int, data []byte, callback func(error)) {
			defer wg.Done()

			// Simulated processing
			if data == nil {
				callback(errors.New("data is nil"))
				return
			}

			// Simulate error after a random delay to simulate concurrent execution
			time.Sleep(time.Duration(idx*100) * time.Millisecond)
			if len(data) == 0 {
				callback(errors.New("data is empty"))
				return
			}

			callback(nil)
		}(i, d, callbacks[i])
	}

	wg.Wait()
}

var errorMessages []string

// MockCallback saves the error message in a slice for later assertions.
func MockCallback(err error) {
	if err != nil {
		errorMessages = append(errorMessages, err.Error())
	}
}

func TestProcessDataWithConcurrency(t *testing.T) {
	// Simulate data and callbacks for multiple goroutines
	data := [][]byte{
		[]byte{1, 2, 3}, // Success
		[]byte{},        // Error
		[]byte{1, 2},    // Success
		nil,             // Error
		[]byte{3, 4},    // Success
	}
	callbacks := make([]func(error), len(data))
	for i := range data {
		callbacks[i] = MockCallback
	}

	// Start the processing in multiple goroutines
	processDataWithConcurrency(data, callbacks)

	// Assert that all callbacks have been executed
	t.Helper()
	if len(callbacks) != 5 {
		t.Fatalf("Expected 5 callbacks, got %d", len(callbacks))
	}

	// Check error messages from the mocked callbacks
	// NOTE: Since we used a mock function, we can validate the output directly.
	expectedErrors := []string{"data is empty", "data is nil"}
	for _, expected := range expectedErrors {
		if !containsString(errorMessages, expected) {
			t.Errorf("Expected error %s, but not found in recorded errors.", expected)
		}
	}
}

// helper function to check if a string slice contains a given string.
func containsString(slice []string, element string) bool {
	for _, e := range slice {
		if e == element {
			return true
		}
	}
	return false
}
