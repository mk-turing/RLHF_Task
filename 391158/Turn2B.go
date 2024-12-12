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
	priority int
	fn       func()
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
			task.fn()
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

// AddTask adds a task to the worker pool with an optional priority
func (pool *WorkerPool) AddTask(priority int, task func()) {
	pool.tasks <- Task{priority: priority, fn: task}
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

	// Simulate tasks with unpredictable execution times and numbers
	for i := 0; i < 100; i++ {
		duration := time.Duration(rand.Intn(500)) * time.Millisecond
		pool.AddTask(rand.Intn(10), func() {
			time.Sleep(duration)
			fmt.Println("Task completed.")
		})
	}

	fmt.Println("All tasks added. Waiting for completion...")
	time.Sleep(5 * time.Second) // Wait for tasks to complete
}
