package main

import (
	"fmt"
	"sync"
	"time"
)

// Task represents a single task that needs to be performed
type Task struct {
	ID        int
	Duration time.Duration
}

// Worker is a goroutine that processes tasks
type Worker struct {
	ID       int
	Tasks    <-chan Task
	wg       *sync.WaitGroup
	metrics  chan<- Task
	done     chan<- struct{}
}

func (w *Worker) Start() {
	for task := range w.Tasks {
		fmt.Printf("Worker %d: Processing task %d\n", w.ID, task.ID)
		// Simulate task processing time
		time.Sleep(task.Duration)
		w.metrics <- task
	}
	w.wg.Done()
	close(w.metrics)
	close(w.done)
}

// LoadBalancer is responsible for distributing tasks across workers
type LoadBalancer struct {
	workers   []*Worker
	nextWorker int
	wg        sync.WaitGroup
}

func (lb *LoadBalancer) AddWorker(worker *Worker) {
	lb.workers = append(lb.workers, worker)
}

func (lb *LoadBalancer) DistributeTask(task Task) {
	worker := lb.workers[lb.nextWorker]
	worker.Tasks <- task
	lb.nextWorker = (lb.nextWorker + 1) % len(lb.workers)