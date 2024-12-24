package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func task(ctx context.Context, wg *sync.WaitGroup, counter *int, workloadFunc func() time.Duration) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Simulate dynamic workload by pausing before each task
			pauseDuration := workloadFunc()
			time.Sleep(pauseDuration)

			for i := 0; i < 1000; i++ {
				*counter++
			}
		}
	}
}

// Workload generator function: Returns a duration with an average of 10ms sleep
func averageWorkloadGenerator() time.Duration {
	return time.Duration(rand.Intn(20)) * time.Millisecond
}

// Workload generator function: Returns a duration with an average of 1ms sleep and occasional bursts of 100ms
func burstyWorkloadGenerator() time.Duration {
	randomNum := rand.Intn(100)
	if randomNum < 10 {
		return time.Duration(100) * time.Millisecond
	}
	return time.Duration(rand.Intn(2)) * time.Millisecond
}

func main() {
	// Number of workers (goroutines)
	numWorkers := 10
	// Total number of tasks each worker will attempt
	totalTasksPerWorker := 1000000
	// Total timeout for all tasks to complete
	totalTimeout := 30 * time.Second
	// Function to generate workload (change to `burstyWorkloadGenerator` for a different scenario)
	workloadFunc := averageWorkloadGenerator

	// Initialize counter to 0
	counter := 0
	wg := sync.WaitGroup{}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), totalTimeout)
	defer cancel()

	// Start timer
	startTime := time.Now()

	// Spawn workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go task(ctx, &wg, &counter, workloadFunc)
	}

	// Wait for all workers to complete or timeout
	wg.Wait()

	// Stop timer
	endTime := time.Now()

	// Calculate metrics
	effectiveTasks := counter
	totalTime := endTime.Sub(startTime)
	throughput := float64(effectiveTasks) / float64(totalTime.Seconds())
	latency := float64(totalTime.Seconds()) / float64(numWorkers*totalTasksPerWorker)

	// Output results
	fmt.Printf("Counter: %d\n", counter)
	fmt.Printf("Total Time (s): %.4f\n", totalTime.Seconds())
	fmt.Printf("Throughput (tasks/s): %.2f\n", throughput)
	fmt.Printf("Approximate Latency (s): %.6f\n", latency)

	if ctx.Err() != nil {
		fmt.Println("Timeout exceeded.")
	}
}
