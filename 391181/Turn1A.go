package main

import (
	"context"
	"fmt"
	"time"
)

// Simulate a long-running task
func longRunningTask(ctx context.Context, id int) {
	select {
	case <-time.After(5 * time.Second): // Simulate work
		fmt.Printf("Task %d completed\n", id)
	case <-ctx.Done():
		fmt.Printf("Task %d canceled: %s\n", id, ctx.Err())
	}
}

func main() {
	// Set a timeout for the entire operation
	timeout := 3 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel() // Ensure resources are cleaned up

	for i := 0; i < 3; i++ {
		// Start a long-running task
		go longRunningTask(ctx, i)
	}

	// Wait for the context to be done
	<-ctx.Done()

	// Optionally, handle the result of the operation
	fmt.Println("All tasks completed or timed out")
}
