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
	fn         func() error
	maxRetries int
	retryDelay time.Duration
	priority   int
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

// worker is a single worker Goroutine that processes tasks
func (pool *WorkerPool) worker() {
	defer pool.wg.Done()
	for {
		select {
		case task, ok := <-pool.tasks:
			if !ok {
				return // Shutdown signal received
			}
			if err := pool.executeTask(task); err != nil {
				fmt.Println("Task failed after all retries:", err)
			}
			atomic.AddInt64(&pool.pendingTasks, -1)
		case <-pool.shutdown:
			return // Shutdown signal received
		}
	}
}

// adjustWorkers dynamically adjusts the number of workers
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
			pool.shutdown <- struct{}{}
		}
	}
}

// executeTask handles retrying the task with a backoff strategy
func (pool *WorkerPool) executeTask(task Task) error {
	for retry := 0; retry <= task.maxRetries; retry++ {
		err := task.fn()
		if err == nil {
			return nil // Task completed successfully
		}

		fmt.Printf("Task failed, retrying (%d/%d): %v\n", retry+1, task.maxRetries+1, err)

		if retry < task.maxRetries {
			time.Sleep(10)
		}
	}
	return fmt.Errorf("max retries reached for task")
}

// AddTask adds a task to the worker pool with retry configuration
func (pool *WorkerPool) AddTask(maxRetries int, retryDelay time.Duration, priority int, task func() error) {
	pool.tasks <- Task{fn: task, maxRetries: maxRetries, retryDelay: retryDelay, priority: priority}
	atomic.AddInt64(&pool.pendingTasks, 1)
}

// Shutdown shuts down the worker pool, waiting for all tasks to complete
func (pool *WorkerPool) Shutdown() {
	close(pool.tasks)
	pool.wg.Wait()
	close(pool.shutdown)
}

func main() {
	// Create a worker pool with dynamic adjustment
	pool := NewWorkerPool(2, 10)
	defer pool.Shutdown()

	// Simulate tasks with unpredictable execution times and occasional failures
	for i := 0; i < 100; i++ {
		duration := time.Duration(rand.Intn(500)) * time.Millisecond
		retries := rand.Intn(3)
		delay := time.Duration(rand.Intn(100)) * time.Millisecond

		pool.AddTask(retries, delay, rand.Intn(10), func() error {
			if rand.Intn(2) == 0 {
				return fmt.Errorf("task failed intentionally")
			}
			time.Sleep(duration)
			fmt.Println("Task completed.")
			return nil
		})
	}

	fmt.Println("All tasks added. Waiting for completion...")
	time.Sleep(5 * time.Second) // Wait for tasks to complete
}
