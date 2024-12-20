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

// Task represents a task to be processed
type Task func()

// WorkerPool manages a fixed number of worker goroutines with dynamic workload balancing
type WorkerPool struct {
	ch     chan Task // Channel to hold tasks
	wg     sync.WaitGroup
	size   int
	ctx    context.Context
	cancel context.CancelFunc
}

// NewWorkerPool creates a new worker pool with a specified number of workers
func NewWorkerPool(size int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		ch:     make(chan Task, size),
		wg:     sync.WaitGroup{},
		size:   size,
		ctx:    ctx,
		cancel: cancel,
	}
}

// Start starts the worker pool
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.size; i++ {
		go wp.worker()
	}
}

// worker is a single worker goroutine that processes tasks from the channel
func (wp *WorkerPool) worker() {
	defer wp.wg.Done()
	for {
		select {
		case task := <-wp.ch:
			task()
		case <-wp.ctx.Done():
			return // Exit early if the context is canceled
		}
	}
}

// Stop stops the worker pool
func (wp *WorkerPool) Stop() {
	close(wp.ch)
	wp.wg.Wait()
	wp.cancel()
}

// Submit submits a task to the worker pool
func (wp *WorkerPool) Submit(task Task) {
	select {
	case wp.ch <- task:
	default:
		// Buffer is full; print a message or take another action
		fmt.Println("Worker pool is currently busy; task will be delayed.")
	}
}

// ProcessLargeDataset processes a large dataset using a worker pool
func ProcessLargeDataset(ctx context.Context, dataset LargeDataset, wp *WorkerPool) {
	chunkSize := len(dataset) / wp.size
	if len(dataset)%wp.size > 0 {
		chunkSize++
	}

	for i := 0; i < len(dataset); i += chunkSize {
		end := i + chunkSize
		if end > len(dataset) {
			end = len(dataset)
		}
		wp.Submit(func() {
			for _, value := range dataset[i:end] {
				select {
				case <-ctx.Done():
					return // Exit early if the context is canceled
				default:
				}
				fmt.Println(value * value) // Example processing operation
			}
		})
	}
}

// Main function to demonstrate concurrent processing of a large dataset using a worker pool
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a large dataset
	largeDataset := make(LargeDataset, 1000000)
	for i := range largeDataset {
		largeDataset[i] = i
	}

	// Create a worker pool with a specified number of workers
	workerPool := NewWorkerPool(runtime.NumCPU())
	workerPool.Start()

	// Process the dataset using the worker pool
	ProcessLargeDataset(ctx, largeDataset, workerPool)

	// Wait for the worker pool to complete all tasks
	workerPool.Stop()
}
