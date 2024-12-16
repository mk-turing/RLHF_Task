package main

import (
	"container/heap"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Task represents a function to be executed by the workers
type Task struct {
	id       int
	priority int // lower values are higher priority
}

// A MinHeap is a min-heap of Tasks.
type MinHeap []*Task

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i].priority < h[j].priority }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MinHeap) Push(x interface{}) {
	*h = append(*h, x.(*Task))
}

func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}

// CircuitBreaker manages the state of the circuit
type CircuitBreaker struct {
	failureThreshold int
	failureCount     int
	isOpen           bool
	resetTimeout     time.Duration
	lastFailureTime  time.Time
	mu               sync.Mutex
}

func NewCircuitBreaker(threshold int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		failureThreshold: threshold,
		resetTimeout:     timeout,
	}
}

func (cb *CircuitBreaker) allowRequest() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	if cb.isOpen {
		if time.Since(cb.lastFailureTime) > cb.resetTimeout {
			cb.isOpen = false // reset the breaker after timeout
		} else {
			return false // still open
		}
	}
	return true
}

func (cb *CircuitBreaker) recordFailure() {
	cb.mu.Lock()
	cb.failureCount++
	cb.lastFailureTime = time.Now()
	if cb.failureCount >= cb.failureThreshold {
		cb.isOpen = true
	}
	cb.mu.Unlock()
}

// Worker represents a worker that processes tasks from a task channel
type Worker struct {
	id             int
	taskChannel    chan *Task
	wg             *sync.WaitGroup
	circuitBreaker *CircuitBreaker
}

// Process the task with retry logic
func (w *Worker) process(task *Task) {
	defer w.wg.Done()
	if !w.circuitBreaker.allowRequest() {
		fmt.Printf("Worker %d: Circuit is open, skipping task %d\n", w.id, task.id)
		return
	}

	attempts := 0
	const maxRetries = 3
	for {
		attempts++
		if err := doWork(task.id); err != nil {
			w.circuitBreaker.recordFailure()
			if attempts >= maxRetries {
				fmt.Printf("Worker %d: Task %d failed after %d attempts: %s\n", w.id, task.id, attempts, err)
				return
			}
			fmt.Printf("Worker %d: Task %d failed, retrying...\n", w.id, task.id)
			continue
		}
		fmt.Printf("Worker %d: Task %d completed successfully.\n", w.id, task.id)
		return
	}
}

// Simulate the actual work, returning an error randomly to mimic failure
func doWork(taskID int) error {
	// 50% chance of failure
	if rand.Intn(2) == 0 {
		return fmt.Errorf("an error occurred with task %d", taskID)
	}
	return nil
}

// WorkerPool manages a pool of workers
type WorkerPool struct {
	taskChannel    chan *Task
	minHeap        MinHeap
	wg             sync.WaitGroup
	numWorkers     int
	circuitBreaker *CircuitBreaker
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(numWorkers int, cb *CircuitBreaker) *WorkerPool {
	pool := &WorkerPool{
		taskChannel:    make(chan *Task),
		numWorkers:     numWorkers,
		circuitBreaker: cb,
		minHeap:        make(MinHeap, 0),
	}
	heap.Init(&pool.minHeap)
	pool.start()
	return pool
}

func (p *WorkerPool) start() {
	for i := 0; i < p.numWorkers; i++ {
		worker := &Worker{id: i, taskChannel: p.taskChannel, wg: &p.wg, circuitBreaker: p.circuitBreaker}
		go func(w *Worker) {
			for {
				// Get the task with the highest priority
				var task *Task
				if len(p.minHeap) > 0 {
					task = heap.Pop(&p.minHeap).(*Task)
				} else {
					// Wait for new tasks, if no tasks are available
					task = <-w.taskChannel
				}
				w.process(task)
			}
		}(worker)
	}
}

// Submit submits a task to the worker pool
func (p *WorkerPool) Submit(task *Task) {
	p.wg.Add(1)
	heap.Push(&p.minHeap, task)
}

// Wait waits for all tasks to finish
func (p *WorkerPool) Wait() {
	p.wg.Wait()
	close(p.taskChannel)
}

// Demonstration of distributed worker pool with shared queue
func main() {
	rand.Seed(time.Now().UnixNano())
	numTasks := 10
	totalWorkers := 3
	circuitBreaker := NewCircuitBreaker(5, 5*time.Second)
	workerPool := NewWorkerPool(totalWorkers, circuitBreaker)

	// Submit tasks to the worker pool with varying priorities
	for i := 0; i < numTasks; i++ {
		priority := rand.Intn(10) // Random priority between 0 (highest) and 9 (lowest)
		workerPool.Submit(&Task{id: i, priority: priority})
	}

	// Wait for all tasks to finish
	workerPool.Wait()
	fmt.Println("All tasks are completed.")
}
