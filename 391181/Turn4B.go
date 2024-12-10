package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Define a transient error type for demonstration
type transientError struct {
	message string
}

func (e transientError) Error() string {
	return e.message
}

// Simulate a long-running task that may return an error
func longRunningTask(ctx context.Context, id int, wg *sync.WaitGroup, errCh chan<- error) {
	defer wg.Done() // Mark this Goroutine as done when the function exits

	maxRetries := 3
	retryDelay := time.Millisecond * 100
	maxRetryDelay := time.Second * 1

	for attempt := 0; attempt <= maxRetries; attempt++ {
		select {
		case <-time.After(time.Duration(id) * time.Second): // Simulate variable work duration
			if id%2 == 0 { // Simulate a transient error for even ID tasks
				if attempt < maxRetries {
					err := transientError{fmt.Sprintf("task %d encountered a transient error, retrying (%d/%d)", id, attempt+1, maxRetries+1)}
					fmt.Println(err)
					// Calculate exponential backoff delay
					delay := time.Duration(rand.Intn(int(maxRetryDelay/retryDelay))) * retryDelay
					time.Sleep(delay)
					continue
				}
				errCh <- fmt.Errorf("task %d failed after max retries: %w", id, transientError{fmt.Sprintf("task %d encountered a transient error", id)})
				return
			}
			fmt.Printf("Task %d completed successfully\n", id)
			return
		case <-ctx.Done():
			fmt.Printf("Task %d canceled: %s\n", id, ctx.Err())
			return
		}
	}
}

func main() {
	// Set a timeout for the entire operation
	timeout := 4 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel() // Ensure that the context is canceled when main exits

	var wg sync.WaitGroup
	errCh := make(chan error, 5) // Buffered channel to collect errors

	for i := 1; i <= 5; i++ {
		wg.Add(1)                              // Increment the WaitGroup counter
		go longRunningTask(ctx, i, &wg, errCh) // Start the long-running task
	}

	// Close the error channel after all tasks are done
	go func() {
		wg.Wait()
		close(errCh)
	}()

	// Wait for all tasks to complete or the context to be done
	<-ctx.Done()

	// Collect and log errors from the error channel
	var aggregatedErrors []error
	for err := range errCh {
		aggregatedErrors = append(aggregatedErrors, err)
	}

	// Log aggregated errors if any
	if len(aggregatedErrors) > 0 {
		fmt.Println("Errors encountered:")
		for _, err := range aggregatedErrors {
			fmt.Println(err)
		}
	} else {
		fmt.Println("All tasks completed successfully.")
	}
}
