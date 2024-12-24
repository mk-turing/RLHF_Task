package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func task(ctx context.Context, wg *sync.WaitGroup, counter *int) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			// If context is canceled, exit the task
			return
		default:
			// Perform the task work
			for i := 0; i < 1000; i++ {
				*counter++
			}
		}
	}
}

func main() {
	// Number of workers (goroutines)
	numWorkers := 10
	// Total number of tasks each worker will attempt
	totalTasksPerWorker := 1000
	// Total timeout for all tasks to complete
	totalTimeout := 5 * time.Second

	// Initialize counter to 0
	counter := 0

	// WaitGroup to wait for all workers to complete or be canceled
	wg := sync.WaitGroup{}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), totalTimeout)
	defer cancel() // Ensure the context is canceled in any case

	// Start timer
	startTime := time.Now()

	// Spawn workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go task(ctx, &wg, &counter)
	}

	// Wait for all workers to complete or timeout
	wg.Wait()

	// Stop timer
	endTime := time.Now()

	// Calculate total time taken
	totalTime := endTime.Sub(startTime)

	// Calculate the effective number of tasks completed
	effectiveTasks := counter

	// Calculate throughput
	throughput := float64(effectiveTasks) / float64(totalTime.Seconds())

	// Calculate latency per task (in seconds)
	latency := float64(totalTime.Seconds()) / float64(numWorkers*totalTasksPerWorker) // Approximate target latency

	// Output results
	fmt.Printf("Counter: %d\n", counter)
	fmt.Printf("Total Time (s): %.4f\n", totalTime.Seconds())
	fmt.Printf("Throughput (tasks/s): %.2f\n", throughput)
	fmt.Printf("Approximate Latency (s): %.6f\n", latency)

	// Check if context was canceled (indicating a timeout)
	if ctx.Err() != nil {
		fmt.Println("Timeout exceeded.")
	}
}
