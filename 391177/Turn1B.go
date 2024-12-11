package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
)

// Task represents a task with a priority and a function to execute
type Task struct {
	Priority int
	Function func()
}

// PriorityQueue manages a collection of tasks by priority
type PriorityQueue struct {
	mu         sync.Mutex
	tasks      []Task
	taskChan   chan Task
	stopChan   chan struct{}
	wg         sync.WaitGroup
	ctx        context.Context
	cancelFunc context.CancelFunc
}

// NewPriorityQueue creates a new PriorityQueue
func NewPriorityQueue() *PriorityQueue {
	ctx, cancel := context.WithCancel(context.Background())
	pq := &PriorityQueue{
		taskChan:   make(chan Task),
		stopChan:   make(chan struct{}),
		wg:         sync.WaitGroup{},
		ctx:        ctx,
		cancelFunc: cancel,
	}
	go pq.worker()
	return pq
}

// Add adds a task to the priority queue
func (pq *PriorityQueue) Add(task Task) {
	pq.taskChan <- task
}

// worker continuously processes tasks from the queue
func (pq *PriorityQueue) worker() {
	pq.wg.Add(1)
	defer pq.wg.Done()
	for {
		select {
		case task, ok := <-pq.taskChan:
			if !ok {
				return
			}
			pq.mu.Lock()
			defer pq.mu.Unlock()
			i := 0
			for ; i < len(pq.tasks); i++ {
				if task.Priority >= pq.tasks[i].Priority {
					break
				}
			}
			pq.tasks = append(pq.tasks[:i], append([]Task{task}, pq.tasks[i:]...)...)
			go pq.executeTask(task)

		case <-pq.stopChan:
			return

		case <-pq.ctx.Done():
			return
		}
	}
}

// executeTask runs the given task
func (pq *PriorityQueue) executeTask(task Task) {
	defer pq.wg.Done()
	task.Function()
}

// Stop stops the priority queue
func (pq *PriorityQueue) Stop() {
	close(pq.taskChan)
	pq.cancelFunc()
	pq.wg.Wait()
	close(pq.stopChan)
}

func main() {
	pq := NewPriorityQueue()

	// Signal handling
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)

	go func() {
		<-signalChan
		fmt.Println("Received signal, stopping priority queue...")
		pq.Stop()
		os.Exit(0)
	}()

	// Add tasks with different priorities
	pq.Add(Task{Priority: 2, Function: func() {
		time.Sleep(2 * time.Second)
		fmt.Println("Task 2 (High Priority) completed")
	}})

	pq.Add(Task{Priority: 1, Function: func() {
		time.Sleep(1 * time.Second)
		fmt.Println("Task 1 (Medium Priority) completed")
	}})

	pq.Add(Task{Priority: 3, Function: func() {
		time.Sleep(3 * time.Second)
		fmt.Println("Task 3 (Low Priority) completed")
	}})

	fmt.Println("Press Ctrl+C to stop")
	select {} // Block the main goroutine
}
