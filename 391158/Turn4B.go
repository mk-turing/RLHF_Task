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
	priority   int // Higher value means higher priority
	attempts   int
	maxRetries int
	fn         func() error
}

// ByPriority is a custom sorter for tasks
type ByPriority []Task

func (t ByPriority) Len() int           { return len(t) }
func (t ByPriority) Less(i, j int) bool { return t[i].priority > t[j].priority }
func (t ByPriority) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t *ByPriority) Pop() (Task, error) {
	if t.Len() == 0 {
		return Task{}, fmt.Errorf("no tasks to pop")
	}
	highestPriorityTask := (*t)[0]
	*t = (*t)[1:] // Remove the first element
	return highestPriorityTask, nil
}

// WorkerPool manages a pool of worker Goroutines
type WorkerPool struct {
	tasks        ByPriority
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
		tasks:      make(ByPriority, 0),
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

// worker processes tasks from the priority queue
func (pool *WorkerPool) worker() {
	defer pool.wg.Done()
	for {
		select {
		case <-pool.shutdown:
			return // Shutdown signal received
		default:
			pool.mu.Lock()
			if len(pool.tasks) == 0 {
				pool.mu.Unlock()
				time.Sleep(10 * time.Millisecond) // Yield if no tasks are available
				continue
			}
			task, _ := pool.tasks.Pop()
			pool.mu.Unlock()

			t := task
			if err := t.fn(); err != nil {
				t.attempts++
				if t.attempts < t.maxRetries {
					go pool.retryTask(t) // Retry in a new Goroutine
				} else {
					fmt.Printf("Task failed after %d attempts: %v\n", t.attempts, err)
				}
			} else {
				atomic.AddInt64(&pool.pendingTasks, -1)
			}
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

// AddTask adds a task to the worker pool with retry capabilities and priority
func (pool *WorkerPool) AddTask(priority int, maxRetries int, taskFunc func() error) {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	//sort.Insert(pool.tasks, task)
	atomic.AddInt64(&pool.pendingTasks, 1)
}

// Shutdown shuts down the worker pool, waiting for all tasks to complete
func (pool *WorkerPool) Shutdown() {
	close(pool.shutdown)
	pool.wg.Wait()
}

func main() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Create a worker pool with dynamic adjustment
	pool := NewWorkerPool(2, 10)
	defer pool.Shutdown()

	// Simulate tasks with unpredictable execution times, failure rates, and priorities
	for i := 0; i < 20; i++ {
		priority := rand.Intn(10) + 1 // Priority 1 to 10
		duration := time.Duration(rand.Intn(500)) * time.Millisecond
		pool.AddTask(priority, 3, func() error { // Maximum 3 retries
			if rand.Intn(2) == 0 { // Simulate random failure
				time.Sleep(duration)
				return fmt.Errorf("simulated error")
			}
			time.Sleep(duration)
			fmt.Printf("Task with priority %d completed successfully.\n", priority)
			return nil
		})
	}

	fmt.Println("All tasks added. Waiting for completion...")
	time.Sleep(10 * time.Second) // Wait for tasks to complete
}
