package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// LargeDataset represents a large dataset
type LargeDataset []int

// WorkerPool manages a fixed number of worker goroutines with dynamic workload balancing
type WorkerPool struct {
	tasks   chan func()      // Channel to hold tasks
	workers chan chan func() // Channel to hold available worker channels
	wg      sync.WaitGroup
	size    int
	ctx     context.Context
	cancel  context.CancelFunc
	balance bool // Enable dynamic workload balancing
}

// NewWorkerPool creates a new worker pool with a specified number of workers
func NewWorkerPool(size int, balance bool) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		tasks:   make(chan func(), size),
		workers: make(chan chan func(), size),
		wg:      sync.WaitGroup{},
		size:    size,
		ctx:     ctx,
		cancel:  cancel,
		balance: balance,
	}
}

// Start starts the worker pool
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.size; i++ {
		wp.wg.Add(1)
		go wp.worker()
	}
}

// worker is a single worker goroutine that processes tasks from the channel
func (wp *WorkerPool) worker() {
	defer wp.wg.Done()
	workerCh := make(chan func())
	wp.workers <- workerCh // Register the worker channel

	for {
		select {
		case task := <-workerCh:
			task()
			wp.workers <- workerCh // Re-register the worker channel for more tasks
		case <-wp.ctx.Done():
			return // Exit early if the context is canceled
		}
	}
}

// Submit submits a task to the worker pool
func (wp *WorkerPool) Submit(task func()) {
	wp.tasks <- task
}

// Dispatch dispatches tasks to available workers
func (wp *WorkerPool) Dispatch() {
	go func() {
		for {
			select {
			case task := <-wp.tasks:
				// Wait for an available worker channel
				workerCh := <-wp.workers
				workerCh <- task // Dispatch the task to the worker
			case <-wp.ctx.Done():
				return // Exit early if the context is canceled
			}
		}
	}()
}

// Stop stops the worker pool
func (wp *WorkerPool) Stop() {
	close(wp.tasks)
	wp.wg.Wait()
	wp.cancel()
}

// ProcessLargeDataset processes a large dataset using a worker pool
func ProcessLargeDataset(ctx context.Context, dataset LargeDataset, wp *WorkerPool) {
	wp.Start()
	wp.Dispatch()

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
	// Wait for all tasks to be processed
	wp.Stop()
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Create a large dataset
	largeDataset := make(LargeDataset, 1000000)
	for i := range largeDataset {
		largeDataset[i] = i
	}
	ProcessLargeDataset(ctx, largeDataset, NewWorkerPool(2, true))
	// Wait for all results to be processed
	<-ctx.Done()
	fmt.Println("Processing completed (or timed out).")
}
