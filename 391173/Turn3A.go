package main

import (
	"context"
	"log"
	"sync"
	"time"
)

// TaskFunc type for task functions
type TaskFunc func(context.Context)

// TaskGroup manages a group of tasks
type TaskGroup struct {
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	tasks      []TaskFunc
	mu         sync.Mutex
	addChan    chan TaskFunc
	removeChan chan int // Index of task to remove
	doneChan   chan struct{}
}

func NewTaskGroup() *TaskGroup {
	ctx, cancel := context.WithCancel(context.Background())
	return &TaskGroup{
		ctx:        ctx,
		cancel:     cancel,
		addChan:    make(chan TaskFunc),
		removeChan: make(chan int),
		doneChan:   make(chan struct{}),
	}
}

func (g *TaskGroup) start() {
	go func() {
		for {
			select {
			case <-g.doneChan:
				return
			case newTask := <-g.addChan:
				g.mu.Lock()
				g.tasks = append(g.tasks, newTask)
				g.wg.Add(1)
				go g.runTask(newTask)
				g.mu.Unlock()
			case index := <-g.removeChan:
				g.mu.Lock()
				if index >= 0 && index < len(g.tasks) {
					tasksLen := len(g.tasks)
					g.tasks[index] = g.tasks[tasksLen-1] // Move the last task to the index
					g.tasks = g.tasks[:tasksLen-1]       // Resize the slice
					log.Printf("Removed task at index %d\n", index)
				}
				g.mu.Unlock()
			}
		}
	}()
}

func (g *TaskGroup) runTask(task TaskFunc) {
	defer g.wg.Done()
	task(g.ctx)
}

func (g *TaskGroup) wait() {
	g.wg.Wait()
	close(g.doneChan)
}

// Add a new task to the group
func (g *TaskGroup) addTask(task TaskFunc) {
	g.addChan <- task
}

// Remove a task by index
func (g *TaskGroup) removeTask(index int) {
	g.removeChan <- index
}

// Stop all tasks
func (g *TaskGroup) stop() {
	g.cancel()
	g.wait()
}

// Example scheduled tasks
func task1(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Task 1: Shutting down")
			return
		case <-ticker.C:
			log.Println("Task 1: Running task")
		}
	}
}

func task2(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Task 2: Shutting down")
			return
		case <-ticker.C:
			log.Println("Task 2: Running task")
		}
	}
}

func main() {
	group := NewTaskGroup()
	group.start()

	// Adding initial tasks
	group.addTask(task1)
	group.addTask(task2)

	// Let tasks run for a few seconds
	time.Sleep(5 * time.Second)

	// Dynamically add a new task
	group.addTask(func(ctx context.Context) {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("Task 3: Shutting down")
				return
			case <-ticker.C:
				log.Println("Task 3: Running task")
			}
		}
	})

	// Let tasks run for a few seconds
	time.Sleep(5 * time.Second)

	// Dynamically remove task at index 0 (task1)
	group.removeTask(0)

	// Let remaining tasks run for a few seconds
	time.Sleep(5 * time.Second)

	// Stop all tasks
	group.stop()

	log.Println("All tasks have been shut down gracefully.")
}
