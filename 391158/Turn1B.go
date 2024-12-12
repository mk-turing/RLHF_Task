package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Task represents a task to be executed by a worker
type Task func()

// WorkerPool manages a pool of worker Goroutines
type WorkerPool struct {
	tasks    chan Task
	workers  int
	wg       sync.WaitGroup
	shutdown chan struct{}
}

// NewWorkerPool creates a new worker pool with the specified number of workers
func NewWorkerPool(workers int) *WorkerPool {
	pool := &WorkerPool{
		tasks:    make(chan Task, workers*10),
		workers:  workers,
		shutdown: make(chan struct{}),
	}

	// Start worker Goroutines
	for i := 0; i < workers; i++ {
		pool.wg.Add(1)
		go pool.worker()
	}

	return pool
}

// worker is a single worker Goroutine that processes tasks
func (pool *WorkerPool) worker() {
	defer pool.wg.Done()
	for {
		select {
		case task, ok := <-pool.tasks:
			if !ok {
				return // Shutdown signal received
			}
			task()
		case <-pool.shutdown:
			return // Shutdown signal received
		}
	}
}

// AddTask adds a task to the worker pool
func (pool *WorkerPool) AddTask(task Task) {
	pool.tasks <- task
}

// Shutdown shuts down the worker pool, waiting for all tasks to complete
func (pool *WorkerPool) Shutdown() {
	close(pool.tasks)
	pool.wg.Wait()
	close(pool.shutdown)
}

func main() {
	// Create a worker pool with 5 workers
	pool := NewWorkerPool(5)
	defer pool.Shutdown()

	// Simulate tasks
	numTasks := 100
	for i := 0; i < numTasks; i++ {
		pool.AddTask(func() {
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			fmt.Println("Task completed.")
		})
	}

	fmt.Println("All tasks added. Waiting for completion...")
}
