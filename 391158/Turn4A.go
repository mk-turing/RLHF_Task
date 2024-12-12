package main

import (
	"container/heap"
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
	priority   int
	fn         func() error
}

// A PriorityQueue implements heap.Interface and holds Tasks.
type PriorityQueue []*Task

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest priority task, so we use greater than here.
	return pq[i].priority > pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*Task))
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	task := old[n-1]
	*pq = old[0 : n-1]
	return task
}

// WorkerPool manages a pool of worker Goroutines
type WorkerPool struct {
	tasks        *PriorityQueue
	tasksLock    sync.Mutex
	minWorkers   int
	maxWorkers   int
	workers      int
	wg           sync.WaitGroup
	shutdown     chan struct{}
	pendingTasks int64
}

// NewWorkerPool creates a new worker pool with dynamic worker adjustment
func NewWorkerPool(minWorkers, maxWorkers int) *WorkerPool {
	pool := &WorkerPool{
		tasks:      &PriorityQueue{},
		minWorkers: minWorkers,
		maxWorkers: maxWorkers,
		shutdown:   make(chan struct{}),
	}

	heap.Init(pool.tasks)

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
		var task *Task
		pool.tasksLock.Lock()
		if pool.tasks.Len() > 0 {
			task = heap.Pop(pool.tasks).(*Task)
		}
		pool.tasksLock.Unlock()

		if task == nil {
			select {
			case <-pool.shutdown:
				return // Shutdown signal received
			default:
				time.Sleep(100 * time.Millisecond) // Sleep to prevent busy waiting
				continue
			}
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
	}
}

// retryTask retries the task and updates the pending tasks counter
func (pool *WorkerPool) retryTask(task *Task) {
	task.fn() // Execute task again
	atomic.AddInt64(&pool.pendingTasks, -1)
}

// adjustWorkers dynamically adjusts the number of workers based on pending tasks
func (pool *WorkerPool) adjustWorkers() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
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
func (pool *WorkerPool) AddTask(priority, maxRetries int, taskFunc func() error) {
	task := &Task{
		fn:         taskFunc,
		maxRetries: maxRetries,
		attempts:   0,
		priority:   priority,
	}
	pool.tasksLock.Lock()
	heap.Push(pool.tasks, task)
	pool.tasksLock.Unlock()
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

	// Simulate tasks with unpredictable execution times and failure rates
	for i := 0; i < 20; i++ {
		duration := time.Duration(rand.Intn(500)) * time.Millisecond
		priority := rand.Intn(10)                // Random priority from 0 to 9
		pool.AddTask(priority, 3, func() error { // Maximum 3 retries
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
