package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// LargeDataset represents a large dataset
type LargeDataset []int

// ProcessLargeDataset processes a large dataset using the range keyword
func ProcessLargeDataset(ctx context.Context, dataset LargeDataset, batchSize int, wg *sync.WaitGroup, results chan int) {
	defer wg.Done()

	for i := 0; i < len(dataset); i += batchSize {
		batch := dataset[i : i+batchSize]
		for _, value := range batch {
			select {
			case <-ctx.Done():
				return // Exit early if the context is canceled
			default:
			}
			results <- value * value // Example processing operation
		}
	}
}

// Main function to demonstrate processing a large dataset
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a large dataset
	largeDataset := make(LargeDataset, 1000000)
	for i := range largeDataset {
		largeDataset[i] = i
	}

	batchSize := 1000
	var wg sync.WaitGroup
	results := make(chan int)

	// Start processing in batches
	numWorkers := 4
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go ProcessLargeDataset(ctx, largeDataset, batchSize, &wg, results)
	}

	// Close the results channel once all workers are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Print processed results
	for result := range results {
		fmt.Println(result)
	}
}
