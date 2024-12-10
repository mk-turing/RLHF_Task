package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Simulate a long-running task
func longRunningTask(ctx context.Context, id int, wg *sync.WaitGroup, errors *[]error) {
	defer wg.Done() // Mark this Goroutine as done when the function exits

	select {
	case <-time.After(time.Duration(id) * time.Second): // Simulate variable work duration
		fmt.Printf("Task %d completed\n", id)
	case <-ctx.Done():
		fmt.Printf("Task %d canceled: %s\n", id, ctx.Err())
		*errors = append(*errors, fmt.Errorf("task %d canceled", id))
	}
}

func main() {
	// Set a timeout for the entire operation
	timeout := 4 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel() // Ensure that the context is canceled when main exits

	var wg sync.WaitGroup
	var errors []error

	for i := 1; i <= 5; i++ {
		wg.Add(1)                                // Increment the WaitGroup counter
		go longRunningTask(ctx, i, &wg, &errors) // Start the long-running task
	}

	// Wait for all tasks to complete or the context to be done
	wg.Wait()

	// Handle and log any accumulated errors
	if len(errors) > 0 {
		fmt.Println("Errors occurred in some tasks:")
		for _, err := range errors {
			fmt.Println(err)
		}
	} else {
		fmt.Println("All tasks completed or timed out without errors")
	}
}
