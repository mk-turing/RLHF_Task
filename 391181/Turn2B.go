package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// Simulate a long-running task that takes a random amount of time
func longRunningTask(ctx context.Context, id int) {
	duration := time.Duration(rand.Intn(5)) * time.Second
	select {
	case <-time.After(duration):
		fmt.Printf("Task %d completed after %s\n", id, duration)
	case <-ctx.Done():
		fmt.Printf("Task %d canceled: %s\n", id, ctx.Err())
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Create a context with a timeout of 4 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel() // Ensure the context is canceled when the function exits

	// Start multiple long-running tasks in Goroutines
	numTasks := 5
	for i := 0; i < numTasks; i++ {
		go longRunningTask(ctx, i)
	}

	// Wait for all tasks to finish or the context to be canceled
	<-ctx.Done()
	fmt.Println("All tasks are complete or timed out.")
}
