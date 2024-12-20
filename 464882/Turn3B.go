package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

// LargeDataset represents a large dataset
type LargeDataset []int

// ProcessLargeDataset processes a large dataset using a worker pool
func ProcessLargeDataset(ctx context.Context, dataset LargeDataset, numWorkers int, results chan int) {
	// Create a worker pool
	workerPool := make(chan int, numWorkers)
	wg := sync.WaitGroup{}
	wg.Add(numWorkers)

	// Spawn workers
	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			for job := range workerPool {
				select {
				case <-ctx.Done():
					return // Exit early if the context is canceled
				default:
				}
				// Example processing operation
				result := job * job
				results <- result
			}
		}()
	}

	// Close the worker pool channel once all jobs are enqueued
	go func() {
		wg.Wait()
		close(workerPool)
	}()

	// Enqueue jobs into the worker pool
	for _, value := range dataset {
		workerPool <- value
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Create a large dataset
	largeDataset := make(LargeDataset, 1000000)
	for i := range largeDataset {
		largeDataset[i] = i
	}
	results := make(chan int)
	numWorkers := 2 * runtime.NumCPU() // Using 2 times the number of CPUs as workers
	go func() {
		for result := range results {
			fmt.Println(result)
		}
	}()
	ProcessLargeDataset(ctx, largeDataset, numWorkers, results)
	// Wait for all results to be processed
	<-ctx.Done()
	fmt.Println("Processing completed (or timed out).")
}
