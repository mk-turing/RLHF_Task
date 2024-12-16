package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Task represents the function to be executed by the workers
type Task struct {
	id int
}

// Worker represents a worker that processes tasks from a task channel
type Worker struct {
	id          int
	taskChannel chan Task
	wg          *sync.WaitGroup
}

// Process the task with retry logic
func (w *Worker) process(task Task) {
	defer w.wg.Done()
	attempts := 0
	const maxRetries = 3

	for {
		attempts++
		// Simulate task processing
		if err := doWork(task.id); err != nil {
			if attempts >= maxRetries {
				fmt.Printf("Worker %d: Task %d failed after %d attempts: %s\n", w.id, task.id, attempts, err)
				return
			}
			// Exponential backoff
			backoffDuration := time.Duration(float64(attempts)*float64(attempts)) * time.Millisecond
			fmt.Printf("Worker %d: Task %d failed, retrying in %v...\n", w.id, task.id, backoffDuration)
			time.Sleep(backoffDuration)
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
	taskChannel chan Task
	wg          sync.WaitGroup
	numWorkers  int
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(numWorkers int) *WorkerPool {
	pool := &WorkerPool{
		taskChannel: make(chan Task),
		numWorkers:  numWorkers,
	}
	pool.start()
	return pool
}

func (p *WorkerPool) start() {
	for i := 0; i < p.numWorkers; i++ {
		worker := &Worker{id: i, taskChannel: p.taskChannel, wg: &p.wg}
		go func(w *Worker) {
			for task := range w.taskChannel {
				w.process(task)
			}
		}(worker)
	}
}

// Submit submits a task to the worker pool
func (p *WorkerPool) Submit(task Task) {
	p.wg.Add(1)
	p.taskChannel <- task
}

// Wait waits for all tasks to finish
func (p *WorkerPool) Wait() {
	p.wg.Wait()
	close(p.taskChannel)
}

// AdjustWorkerPool allows you to dynamically adjust the number of workers
func (p *WorkerPool) AdjustWorkerPool(newSize int) {
	if newSize > p.numWorkers {
		for i := p.numWorkers; i < newSize; i++ {
			worker := &Worker{id: i, taskChannel: p.taskChannel, wg: &p.wg}
			go func(w *Worker) {
				for task := range w.taskChannel {
					w.process(task)
				}
			}(worker)
		}
	}
	p.numWorkers = newSize
}

func main() {
	rand.Seed(time.Now().UnixNano())
	numTasks := 10
	totalWorkers := 3
	workerPool := NewWorkerPool(totalWorkers)

	// Submit tasks to the worker pool
	for i := 0; i < numTasks; i++ {
		workerPool.Submit(Task{id: i})
	}

	// You can dynamically adjust the workers later if needed
	time.Sleep(5 * time.Second)    // Simulate work
	workerPool.AdjustWorkerPool(6) // Increase the worker pool

	// Wait for all tasks to finish
	workerPool.Wait()
	fmt.Println("All tasks are completed.")
}
