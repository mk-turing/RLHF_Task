package main

import (
	"context"
	"log"
	"sync"
	"time"
)

// Task type to define different scheduled tasks
type Task struct {
	ID       int
	Interval time.Duration
}

func runTask(ctx context.Context, task Task, wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(task.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("Task %d: Shutting down\n", task.ID)
			return
		case t := <-ticker.C:
			log.Printf("Task %d: Executing at %s\n", task.ID, t)
		}
	}
}

func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Define multiple tasks with different intervals
	tasks := []Task{
		{ID: 1, Interval: 1 * time.Second},
		{ID: 2, Interval: 2 * time.Second},
		{ID: 3, Interval: 3 * time.Second},
	}

	// Start all scheduled tasks
	for _, task := range tasks {
		wg.Add(1)
		go runTask(ctx, task, &wg)
	}

	// Simulate application running for a duration
	time.Sleep(10 * time.Second) // Let tasks run for a while
	log.Println("Initiating graceful shutdown...")

	// Cancel the context to shut down tasks
	cancel()

	// Wait for all tasks to complete before exiting
	wg.Wait()
	log.Println("All tasks completed. Application exited gracefully.")
}
