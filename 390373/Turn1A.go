package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"
)

// Custom error structure to hold retry metadata
type RetryableError struct {
	Err        error
	RetryCount int
}

// Implement the Error method for RetryableError
func (re *RetryableError) Error() string {
	return fmt.Sprintf("error: %v, retry count: %d", re.Err, re.RetryCount)
}

// Example function simulating an operation that may fail
func operation(ctx context.Context) error {
	// Randomly simulate a failure for demonstration
	if rand.Float32() < 0.7 { // ~70% chance to fail
		return &RetryableError{
			Err:        fmt.Errorf("operation failed"),
			RetryCount: 0, // This value will be updated externally
		}
	}
	return nil
}

// Retry logic with observable error handling
func retryOperation(ctx context.Context, maxRetries int, delay time.Duration) error {
	retries := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := operation(ctx)
			if err == nil {
				return nil // success
			}

			// Use type assertion to extract the RetryableError
			if retryableErr, ok := err.(*RetryableError); ok {
				retryableErr.RetryCount = retries // update count
				log.Printf("Attempting to retry (%d) due to error: %v", retries, retryableErr)
				retries++
				if retries > maxRetries {
					return fmt.Errorf("max retries exceeded: %w", retryableErr)
				}
				time.Sleep(delay) // wait before retrying
				continue          // retry the operation
			}
			return fmt.Errorf("non-retryable error occurred: %w", err)
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano()) // Seed random number for simulating operation's success/failure

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Context with timeout
	defer cancel()

	err := retryOperation(ctx, 5, 2*time.Second) // Retry up to 5 times with a 2 second delay
	if err != nil {
		log.Fatalf("Operation failed: %v", err)
	}

	log.Println("Operation completed successfully")
}
