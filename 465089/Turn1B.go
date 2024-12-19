package main

import (
	"context"
	"fmt"
	"time"
)

// Define a custom error type for retryable errors
type RetryableError struct {
	message string
}

func (e *RetryableError) Error() string {
	return e.message
}

// Sample function that can fail and requires retry
func sampleFunction(ctx context.Context, input string) (string, error) {
	// Simulate some business logic that might fail
	if input == "retry" {
		return "", &RetryableError{"Sample function failed due to retryable error"}
	}
	return "success", nil
}

// Retry function with context and timeouts
func retry(ctx context.Context, attempt int, f func(context.Context, string) (string, error), input string) (string, error) {
	var result string
	var err error

	// Loop for the specified number of attempts
	for i := 0; i < attempt; i++ {
		// Check if the context has timed out before making the call
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
			// Make the function call
			result, err = f(ctx, input)
			if err == nil {
				// No error, return the result
				return result, nil
			}
			// Handle retryable errors
			if _, ok := err.(*RetryableError); ok {
				// Log the error or do any necessary cleanup before retrying
				fmt.Println("Retryable error:", err.Error())
				// Sleep for a while before retrying
				time.Sleep(time.Second * time.Duration(i))
			} else {
				// Non-retryable error, return the error
				return "", err
			}
		}
	}
	// Exhausted all attempts, return the last error
	return "", err
}

// Example usage of retry function
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Sample input data
	input := "retry"

	// Retry the sampleFunction with 3 attempts and context
	result, err := retry(ctx, 3, sampleFunction, input)
	if err != nil {
		// Handle the final error
		fmt.Println("Final error:", err.Error())
	} else {
		// Handle the successful result
		fmt.Println("Success:", result)
	}
}
