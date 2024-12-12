package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// Task represents a task to be executed by a worker
type Task struct {
	attempts   int
	maxRetries int
	fn         func() error
}

// WorkerPool manages a pool of worker Goroutines
type WorkerPool struct {
	tasks        chan Task
	minWorkers   int
	maxWorkers   int
	workers      int
	wg           sync.WaitGroup
	shutdown     chan struct{}
	mu           sync.Mutex
	pendingTasks int64
}

// NewWorkerPool creates a new worker pool with dynamic worker adjustment
func NewWorkerPool(minWorkers, maxWorkers int) *WorkerPool {
	pool := &WorkerPool{
		tasks:      make(chan Task),
		minWorkers: minWorkers,
		maxWorkers: maxWorkers,
		workers:    minWorkers,
		shutdown:   make(chan struct{}),
	}

	// Start worker Goroutines
	for i := 0; i < minWorkers; i++ {
		pool.wg.Add(1)
		go pool.worker()
	}

	go pool.adjustWorkers()

	return pool
}

// worker processes tasks from the channel
func (pool *WorkerPool) worker() {
	defer pool.wg.Done()
	for {
		select {
		case task, ok := <-pool.tasks:
			if !ok {
				return // Shutdown signal received
			}
			if err := task.fn(); err != nil {
				task.attempts++
				if task.attempts < task.maxRetries {
					go pool.retryTask(task) // Retry in a new Goroutine
				} else {
					fmt.Printf("Task failed after %d attempts: %v\n", task.attempts, err)
				}
			} else {
				atomic.AddInt64(&pool.pendingTasks, -1)
			}
		case <-pool.shutdown:
			return // Shutdown signal received
		}
	}
}

// retryTask retries the task and updates the pending tasks counter
func (pool *WorkerPool) retryTask(task Task) {
	task.fn() // Execute task again
	atomic.AddInt64(&pool.pendingTasks, -1)
}

// adjustWorkers dynamically adjusts the number of workers based on pending tasks
func (pool *WorkerPool) adjustWorkers() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		pool.mu.Lock()
		defer pool.mu.Unlock()

		pending := atomic.LoadInt64(&pool.pendingTasks)

		if pending > int64(pool.workers*2) && pool.workers < pool.maxWorkers {
			fmt.Println("Adding worker:", pool.workers)
			pool.workers++
			pool.wg.Add(1)
			go pool.worker()
		} else if pending < int64(pool.workers/2) && pool.workers > pool.minWorkers {
			fmt.Println("Removing worker:", pool.workers)
			pool.workers--
		}
	}
}

// AddTask adds a task to the worker pool with retry capabilities
func (pool *WorkerPool) AddTask(maxRetries int, taskFunc func() error) {
	pool.tasks <- Task{
		fn:         taskFunc,
		maxRetries: maxRetries,
		attempts:   0,
	}
	atomic.AddInt64(&pool.pendingTasks, 1)
}

// Shutdown shuts down the worker pool, waiting for all tasks to complete
func (pool *WorkerPool) Shutdown() {
	close(pool.tasks)
	pool.wg.Wait()
	close(pool.shutdown)
}

func main() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Create a worker pool with dynamic adjustment
	pool := NewWorkerPool(2, 10)
	defer pool.Shutdown()

	// Simulate tasks with unpredictable execution times and failure rates
	for i := 0; i < 20; i++ {
		duration := time.Duration(rand.Intn(500)) * time.Millisecond
		pool.AddTask(3, func() error { // Maximum 3 retries
			if rand.Intn(2) == 0 { // Simulate random failure
				time.Sleep(duration)
				return fmt.Errorf("simulated error")
			}
			time.Sleep(duration)
			fmt.Println("Task completed successfully.")
			return nil
		})
	}

	fmt.Println("All tasks added. Waiting for completion...")
	time.Sleep(10 * time.Second) // Wait for tasks to complete
}
