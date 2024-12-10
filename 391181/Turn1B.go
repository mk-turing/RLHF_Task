package main

import (
	"context"
	"log"
	"time"
)

// Simulate a long-running operation
func longRunningOperation(ctx context.Context, id int) {
	log.Printf("Operation %d started\n", id)
	select {
	case <-time.After(5 * time.Second): // Simulate work that takes 5 seconds
		log.Printf("Operation %d completed\n", id)
	case <-ctx.Done():
		log.Printf("Operation %d cancelled\n", id)
		return
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // Set a 3-second timeout
	defer cancel()

	// Start two long-running operations in Goroutines
	go longRunningOperation(ctx, 1)
	go longRunningOperation(ctx, 2)

	time.Sleep(2 * time.Second) // Wait for 2 seconds

	// Cancel the context to stop the operations
	log.Println("Cancelling operations...")
	cancel()

	// Wait for Goroutines to finish
	time.Sleep(1 * time.Second)
}
