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

// ProcessLargeDataset processes a large dataset concurrently using the range keyword
func ProcessLargeDataset(ctx context.Context, dataset LargeDataset, wg *sync.WaitGroup) {
	defer wg.Done()

	for _, value := range dataset {
		select {
		case <-ctx.Done():
			return // Exit early if the context is canceled
		default:
		}
		fmt.Println(value * value) // Example processing operation
	}
}

// Main function to demonstrate concurrent processing of a large dataset
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a large dataset
	largeDataset := make(LargeDataset, 1000000)
	for i := range largeDataset {
		largeDataset[i] = i
	}

	var wg sync.WaitGroup

	// Determine the number of workers based on the number of available CPU cores
	numWorkers := runtime.NumCPU()

	// Divide the dataset evenly among workers
	chunkSize := len(largeDataset) / numWorkers
	if len(largeDataset)%numWorkers > 0 {
		chunkSize++
	}

	// Start processing in chunks concurrently
	for i := 0; i < len(largeDataset); i += chunkSize {
		end := i + chunkSize
		if end > len(largeDataset) {
			end = len(largeDataset)
		}
		wg.Add(1)
		go ProcessLargeDataset(ctx, largeDataset[i:end], &wg)
	}

	// Wait for all workers to complete
	wg.Wait()
}
