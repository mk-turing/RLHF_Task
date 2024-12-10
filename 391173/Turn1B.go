package main

import (
	"context"
	"log"
	"sync"
	"time"
)

func scheduledTask(ctx context.Context, taskID int, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("Task %d: Shutting down\n", taskID)
			return
		case <-ticker.C:
			log.Printf("Task %d: Running task at %s\n", taskID, time.Now())
		}
	}
}

func main() {
	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start some scheduled tasks
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		scheduledTask(ctx, 1, time.Second)
		wg.Done()
	}()

	go func() {
		scheduledTask(ctx, 2, 2*time.Second)
		wg.Done()
	}()

	// Simulate a graceful shutdown
	time.Sleep(5 * time.Second) // Let tasks run for a few seconds
	log.Println("Initiating graceful shutdown...")
	cancel() // Cancel the context, signalling goroutines to stop

	// Wait for all tasks to complete
	wg.Wait()
	log.Println("Application exited gracefully.")
}
