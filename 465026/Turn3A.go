package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func task(ctx context.Context, wg *sync.WaitGroup, counter *int) {
	defer wg.Done()

	for i := 0; i < 1000; i++ {
		// Simulate work
		for j := 0; j < 1000; j++ {
			*counter++
		}
		// Randomize task duration
		time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
	}
}

func main() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// List of varying workloads
	workloads := []struct {
		numWorkers   int
		totalTasks   int
		totalTimeout time.Duration
	}{
		{1, 10000, 5 * time.Second},
		{5, 10000, 5 * time.Second},
		{10, 10000, 5 * time.Second},
		{20, 20000, 10 * time.Second},
		{50, 20000, 15 * time.Second},
	}

	for _, workload := range workloads {
		// Initialize counter to 0
		counter := 0

		// WaitGroup to wait for all workers to complete or be canceled
		wg := sync.WaitGroup{}

		// Create a context with a timeout
		ctx, cancel := context.WithTimeout(context.Background(), workload.totalTimeout)
		defer cancel() // Ensure the context is canceled in any case

		// Start timer
		startTime := time.Now()

		// Spawn workers
		for i := 0; i < workload.numWorkers; i++ {
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
		latency := float64(totalTime.Seconds()) / float64(workload.numWorkers*workload.totalTasks)

		// Output results for the current workload
		fmt.Printf("\nWorkload:\n")
		fmt.Printf("Number of Workers: %d\n", workload.numWorkers)
		fmt.Printf("Total Tasks: %d\n", workload.totalTasks)
		fmt.Printf("Total Timeout: %s\n", workload.totalTimeout)
		fmt.Printf("Counter: %d\n", counter)
		fmt.Printf("Total Time (s): %.4f\n", totalTime.Seconds())
		fmt.Printf("Throughput (tasks/s): %.2f\n", throughput)
		fmt.Printf("Latency (s): %.6f\n", latency)

		// Check if context was canceled (indicating a timeout)
		if ctx.Err() != nil {
			fmt.Println("Timeout exceeded.")
		}
	}
}
