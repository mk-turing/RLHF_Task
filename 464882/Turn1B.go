package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ProcessData represents a function that processes a single data item.
type ProcessData func(int)

// processDataChunk processes a chunk of data concurrently using goroutines.
func processDataChunk(ctx context.Context, wg *sync.WaitGroup, data []int, process ProcessData) {
	defer wg.Done()
	for _, item := range data {
		select {
		case <-ctx.Done():
			return // Return if the context is canceled
		default:
			process(item)
		}
	}
}

// processLargeDataset processes a large dataset concurrently using channel buffering.
func processLargeDataset(ctx context.Context, data []int, process ProcessData, chunkSize int) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	// Create a channel with a buffer to hold chunks of data
	dataCh := make(chan []int, 100)

	go func() {
		defer close(dataCh)
		for i := 0; i < len(data); i += chunkSize {
			end := i + chunkSize
			if end > len(data) {
				end = len(data)
			}
			dataCh <- data[i:end]
		}
	}()

	// Start a fixed number of goroutines to process chunks concurrently
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for chunk := range dataCh {
				processDataChunk(ctx, wg, chunk, process)
			}
		}()
	}
}

// main function for demonstration purposes.
func main() {
	// Generate a large dataset for testing
	data := make([]int, 1_000_000)
	for i := range data {
		data[i] = i
	}

	// Define a simple process data function that prints the item
	processData := func(item int) {
		fmt.Println(item)
	}

	// Create a context with a timeout to simulate a long-running task
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Process the large dataset using the processLargeDataset function
	processLargeDataset(ctx, data, processData, 10_000)
}
