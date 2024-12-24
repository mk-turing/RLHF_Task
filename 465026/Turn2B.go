package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Function to represent a task
func task(ctx context.Context, wg *sync.WaitGroup, counter *int) {
	defer wg.Done()
	for i := 0; i < 1000; i++ {
		select {
		case <-ctx.Done():
			return // Exit early if context is canceled
		default:
			*counter++
		}
	}
}

func main() {
	// Number of workers (goroutines)
	numWorkers := 10
	// Total number of tasks each worker will perform
	totalTasks := 100000
	// Timeout duration for each task
	taskTimeout := 10 * time.Millisecond

	// Initialize counter to 0
	counter := 0

	// WaitGroup to wait for all workers to complete
	wg := sync.WaitGroup{}

	// Create a new context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), taskTimeout)
	defer cancel()

	// Start timer
	startTime := time.Now()

	// Spawn workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go task(ctx, &wg, &counter)
	}

	// Wait for all workers to complete (or timeout)
	wg.Wait()

	// Stop timer
	endTime := time.Now()

	// Calculate total time taken
	totalTime := endTime.Sub(startTime)

	// Calculate throughput
	throughput := float64(totalTasks) / float64(totalTime.Seconds())

	// Calculate latency per task (in seconds)
	latency := float64(totalTime.Seconds()) / float64(numWorkers*1000)

	// Check if any task timed out
	timedOut := ctx.Err() == context.DeadlineExceeded

	// Output results
	fmt.Printf("Counter: %d\n", counter)
	fmt.Printf("Total Time (s): %.4f\n", totalTime.Seconds())
	fmt.Printf("Throughput (tasks/s): %.2f\n", throughput)
	fmt.Printf("Latency (s): %.6f\n", latency)
	if timedOut {
		fmt.Println("Note: Some tasks timed out.")
	}
}
