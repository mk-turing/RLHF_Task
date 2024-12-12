package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Task represents a task to be executed by a worker
type Task func()

// WorkerPool manages a dynamic pool of worker Goroutines
type WorkerPool struct {
	tasks       chan Task
	workers     sync.WaitGroup
	workerCount int
	shutdown    chan struct{}
}

const (
	minWorkers = 1
	maxWorkers = 10
)

// NewWorkerPool creates a new dynamic worker pool
func NewWorkerPool() *WorkerPool {
	pool := &WorkerPool{
		tasks:    make(chan Task),
		shutdown: make(chan struct{}),
	}

	// Start a fixed number of workers initially
	for i := 0; i < minWorkers; i++ {
		pool.addWorker()
	}

	go pool.monitor()

	return pool
}

// addWorker adds a new worker to the pool
func (pool *WorkerPool) addWorker() {
	pool.workerCount++
	pool.workers.Add(1)
	go func() {
		defer pool.workers.Done()
		for {
			select {
			case task, ok := <-pool.tasks:
				if !ok {
					return // Channel closed
				}
				task()
			case <-pool.shutdown:
				return // Shutdown signal received
			}
		}
	}()
}

// monitor checks the number of tasks and adds or removes workers based on load
func (pool *WorkerPool) monitor() {
	for {
		select {
		case <-time.After(1 * time.Second):
			taskCount := len(pool.tasks)
			if taskCount > pool.workerCount && pool.workerCount < maxWorkers {
				pool.addWorker()
			} else if taskCount < pool.workerCount/2 && pool.workerCount > minWorkers {
				pool.workerCount-- // Reduce the count; we do not remove the worker yet
				go func() {        // Wait for floating workers to finish
					time.Sleep(2 * time.Second) // Assume some idle time before removing
					if pool.workerCount > minWorkers {
						pool.workerCount--
					}
				}()
			}
		case <-pool.shutdown:
			return
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
	pool.workers.Wait()
	close(pool.shutdown)
}

// Simulate random task duration
func createRandomTask(id int) Task {
	return func() {
		duration := time.Duration(rand.Intn(1000)) * time.Millisecond
		time.Sleep(duration)
		fmt.Printf("Task %d completed in %v\n", id, duration)
	}
}

func main() {
	// Create a dynamic worker pool
	pool := NewWorkerPool()
	defer pool.Shutdown()

	// Simulate adding a burst of tasks
	for i := 0; i < 50; i++ {
		pool.AddTask(createRandomTask(i))
	}

	// Wait a bit to let all tasks complete
	time.Sleep(10 * time.Second)
	fmt.Println("Finished submitting tasks.")
}
