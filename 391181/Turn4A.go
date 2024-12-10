package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Simulate a long-running task that may return transient errors
func longRunningTask(ctx context.Context, id int, wg *sync.WaitGroup, errCh chan<- error) {
	defer wg.Done() // Mark this Goroutine as done when the function exits

	const maxRetries = 3
	baseDelay := 500 * time.Millisecond // Starting delay for retries

	var err error
	for attempt := 0; attempt < maxRetries; attempt++ {
		select {
		case <-ctx.Done():
			fmt.Printf("Task %d canceled: %s\n", id, ctx.Err())
			return
		default:
			// Simulate a task that fails randomly
			err = simulateTask(id)
			if err == nil {
				fmt.Printf("Task %d completed successfully\n", id)
				return
			}

			// Log the error and apply exponential backoff
			fmt.Printf("Task %d encountered an error (attempt %d): %s\n", id, attempt+1, err)

			// Calculate the backoff duration
			backoffDuration := baseDelay * (1 << attempt) // Exponential backoff
			select {
			case <-time.After(backoffDuration):
				// Wait for the backoff duration
			case <-ctx.Done():
				fmt.Printf("Task %d canceled during backoff: %s\n", id, ctx.Err())
				return
			}
		}
	}

	// If all retries failed, send the error to the error channel
	errCh <- fmt.Errorf("task %d failed after %d attempts: %w", id, maxRetries, err)
}

// Simulate a task that randomly fails (transient error)
func simulateTask(id int) error {
	if rand.Float32() < 0.5 { // 50% chance of transient error
		return fmt.Errorf("transient error in task %d", id)
	}
	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator

	// Set a timeout for the entire operation
	timeout := 10 * time.Second
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
