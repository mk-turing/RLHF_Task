package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// LargeDataset represents a large dataset
type LargeDataset []int

// ProcessBatch processes a single batch of the dataset
func ProcessBatch(ctx context.Context, batch LargeDataset, results chan int) {
	for _, value := range batch {
		select {
		case <-ctx.Done():
			return // Exit early if the context is canceled
		default:
		}
		results <- value * value // Example processing operation
	}
}

// ProcessLargeDataset processes a large dataset using concurrent batches
func ProcessLargeDataset(ctx context.Context, dataset LargeDataset, batchSize int, wg *sync.WaitGroup, results chan int) {
	defer wg.Done()

	numBatches := len(dataset) / batchSize
	for i := 0; i < numBatches; i++ {
		batch := dataset[i*batchSize : (i+1)*batchSize]
		go ProcessBatch(ctx, batch, results)
	}
}

// Main function to demonstrate processing a large dataset concurrently
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
