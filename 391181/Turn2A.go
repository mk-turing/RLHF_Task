package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Simulate a long-running task
func longRunningTask(ctx context.Context, id int, wg *sync.WaitGroup) {
	defer wg.Done() // Mark this Goroutine as done when the function exits

	select {
	case <-time.After(time.Duration(id) * time.Second): // Simulate variable work duration
		fmt.Printf("Task %d completed\n", id)
	case <-ctx.Done():
		fmt.Printf("Task %d canceled: %s\n", id, ctx.Err())
	}
}

func main() {
	// Set a timeout for the entire operation
	timeout := 4 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel() // Ensure that the context is canceled when main exits

	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)                       // Increment the WaitGroup counter
		go longRunningTask(ctx, i, &wg) // Start the long-running task
	}

	// Wait for all tasks to complete or the context to be done
	wg.Wait()

	// Notify that the main function is done
	fmt.Println("All tasks completed or timed out")
}
