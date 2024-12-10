package main

import (
	"container/heap"
	"context"
	"log"
	"sync"
	"time"
)

// Task represents a scheduled task with priority and status
type Task struct {
	id       int
	priority int
	status   string
	work     func(context.Context)
	index    int // The index in the priority queue
}

// TaskQueue implements a priority queue for tasks
type TaskQueue []*Task

func (pq TaskQueue) Len() int { return len(pq) }

func (pq TaskQueue) Less(i, j int) bool {
	// Higher priority comes first (lower number means higher priority)
	return pq[i].priority < pq[j].priority
}

func (pq TaskQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *TaskQueue) Push(x interface{}) {
	n := len(*pq)
	task := x.(*Task)
	task.index = n
	*pq = append(*pq, task)
}

func (pq *TaskQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	task := old[n-1]
	old[n-1] = nil  // Avoid memory leak
	task.index = -1 // Mark as removed
	*pq = old[0 : n-1]
	return task
}

// TaskManager manages a group of tasks
type TaskManager struct {
	ctx        context.Context
	cancel     context.CancelFunc
	tasks      TaskQueue
	mu         sync.Mutex
	wg         sync.WaitGroup
	addChan    chan *Task
	removeChan chan int
	doneChan   chan struct{}
}

func NewTaskManager() *TaskManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &TaskManager{
		ctx:        ctx,
		cancel:     cancel,
		addChan:    make(chan *Task),
		removeChan: make(chan int),
		doneChan:   make(chan struct{}),
	}
}

func (tm *TaskManager) start() {
	heap.Init(&tm.tasks)
	go func() {
		for {
			select {
			case <-tm.doneChan:
				return
			case newTask := <-tm.addChan:
				tm.mu.Lock()
				heap.Push(&tm.tasks, newTask)
				tm.mu.Unlock()
				tm.wg.Add(1)
				go tm.runTask(newTask)
			case index := <-tm.removeChan:
				tm.mu.Lock()
				if index >= 0 && index < len(tm.tasks) {
					// Mark task as removed
					task := tm.tasks[index]
					task.status = "removed"
					log.Printf("Task ID %d removed\n", task.id)
				}
				tm.mu.Unlock()
			}
		}
	}()
}

func (tm *TaskManager) runTask(task *Task) {
	defer tm.wg.Done()
	task.status = "running"
	defer func() {
		task.status = "completed"
	}()
	task.work(tm.ctx)
}

func (tm *TaskManager) wait() {
	tm.wg.Wait()
	close(tm.doneChan)
}

// Add a new task to the manager
func (tm *TaskManager) addTask(task *Task) {
	tm.addChan <- task
}

// Remove a task by index
func (tm *TaskManager) removeTask(index int) {
	tm.removeChan <- index
}

// Stop all tasks, prioritizing based on conditions
func (tm *TaskManager) stop() {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	log.Println("Initiating shutdown...")

	// Sort tasks by priority for orderly shutdown
	for tm.tasks.Len() > 0 {
		task := heap.Pop(&tm.tasks).(*Task)
		log.Printf("Stopping task ID %d with priority %d\n", task.id, task.priority)
		task.status = "stopped"
	}
	tm.cancel()
	tm.wait()
	log.Println("All tasks have been shut down gracefully.")
}

// Example task functions
func exampleTask(id int, duration time.Duration) func(context.Context) {
	return func(ctx context.Context) {
		ticker := time.NewTicker(duration)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Printf("Task ID %d: Shutting down\n", id)
				return
			case <-ticker.C:
				log.Printf("Task ID %d: Running\n", id)
			}
		}
	}
}

func main() {
	manager := NewTaskManager()
	manager.start()

	// Adding initial tasks with different priorities
	manager.addTask(&Task{id: 1, priority: 2, work: exampleTask(1, 1*time.Second)})
	manager.addTask(&Task{id: 2, priority: 1, work: exampleTask(2, 2*time.Second)})
	manager.addTask(&Task{id: 3, priority: 3, work: exampleTask(3, 3*time.Second)})

	// Let tasks run for a few seconds
	time.Sleep(5 * time.Second)

	// Dynamically remove task ID 1
	manager.removeTask(0) // This will mark task 1 as removed

	// Let remaining tasks run for a few seconds
	time.Sleep(5 * time.Second)

	// Stop all tasks
	manager.stop()
}
