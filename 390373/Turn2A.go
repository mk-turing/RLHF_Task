package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"
)

// Custom error structure for retryable errors
type RetryableError struct {
	Err        error
	RetryCount int
}

// Implement the Error method for RetryableError
func (re *RetryableError) Error() string {
	return fmt.Sprintf("error: %v, retry count: %d", re.Err, re.RetryCount)
}

// Example error categorization structure
type NonRetryableError struct {
	Err error
}

func (nre *NonRetryableError) Error() string {
	return fmt.Sprintf("non-retryable error: %v", nre.Err)
}

// Simulates an operation that may fail
func operation(ctx context.Context) error {
	if rand.Float32() < 0.7 { // ~70% chance to fail
		// Randomly return a type of error
		if rand.Float32() < 0.5 {
			return &RetryableError{
				Err: fmt.Errorf("temporary error occurred"),
			}
		}
		return &NonRetryableError{
			Err: fmt.Errorf("permanent error occurred"),
		}
	}
	return nil // operation succeeded
}

// Exponential Backoff Retry Logic
func retryOperation(ctx context.Context, maxRetries int) error {
	retries := 0
	backoff := time.Second // starting backoff duration

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := operation(ctx)
			if err == nil {
				return nil // success
			}

			// Handle different error types
			if retryableErr, ok := err.(*RetryableError); ok {
				retryableErr.RetryCount = retries // update count
				log.Printf("Retryable error occurred: %s", retryableErr)
				log.Printf("Attempting retry %d after %v", retries+1, backoff)
				retries++
				if retries > maxRetries {
					return fmt.Errorf("max retries exceeded: %w", retryableErr)
				}
				time.Sleep(backoff) // wait before retrying

				// Exponential backoff
				backoff *= 2 // double the backoff time
				continue     // retry the operation
			}

			// If it's a non-retryable error
			if _, ok := err.(*NonRetryableError); ok {
				log.Printf("Non-retryable error encountered: %s", err)
				return fmt.Errorf("operation failed with non-retryable error: %w", err)
			}

			// Log unexpected error (should not happen)
			log.Printf("Unexpected error: %s", err)
			return fmt.Errorf("unexpected error occurred: %w", err)
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano()) // Seed random number generator

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Context with timeout
	defer cancel()

	err := retryOperation(ctx, 5) // Retry up to 5 times
	if err != nil {
		log.Fatalf("Operation failed: %v", err)
	}

	log.Println("Operation completed successfully")
}
