package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const taskDuration = 5 * time.Second // Duration of each task

// Task represents a scheduled task.
type Task struct {
	ID int
}

// performTask simulates task execution
func performTask(ctx context.Context, task Task) {
	fmt.Printf("Starting task %d\n", task.ID)
	select {
	case <-time.After(taskDuration): // Simulate task work
		fmt.Printf("Completed task %d\n", task.ID)
	case <-ctx.Done(): // Handle shutdown
		fmt.Printf("Task %d cancelled\n", task.ID)
	}
}

func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	// Handle shutdown signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Start Goroutines for scheduled tasks
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(task Task) {
			defer wg.Done()
			performTask(ctx, task)
		}(Task{ID: i})
	}

	// Wait for a shutdown signal
	<-signalChan
	fmt.Println("\nShutdown signal received. Initiating graceful shutdown...")

	// Cancel context to signal Goroutines to stop
	cancel()

	// Wait for all Goroutines to finish
	wg.Wait()
	fmt.Println("All tasks completed. Exiting.")
}
