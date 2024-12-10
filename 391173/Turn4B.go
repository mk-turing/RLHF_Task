package main

import (
	"context"
	"log"
	"sync"
	"time"
)

// TaskFunc type for task functions
type TaskFunc func(context.Context, string)

// Task represents a scheduled task with metadata
type Task struct {
	id       int
	name     string
	priority int // Higher number means higher priority
	ctx      context.Context
	cancel   context.CancelFunc
	done     chan struct{}
}

// TaskGroup manages a group of tasks
type TaskGroup struct {
	mu           sync.Mutex
	tasks        []*Task
	addChan      chan *Task
	removeChan   chan int // Index of task to remove
	shutdownChan chan struct{}
}

func NewTaskGroup() *TaskGroup {
	return &TaskGroup{
		addChan:      make(chan *Task),
		removeChan:   make(chan int),
		shutdownChan: make(chan struct{}),
	}
}

func (g *TaskGroup) start() {
	go func() {
		for {
			select {
			case newTask := <-g.addChan:
				g.mu.Lock()
				g.tasks = append(g.tasks, newTask)
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
			case <-g.shutdownChan:
				g.mu.Lock()
				for _, task := range g.tasks {
					task.cancel()
				}
				g.tasks = nil
				g.mu.Unlock()
				return
			}
		}
	}()
}

func (g *TaskGroup) runTask(task *Task) {
	defer close(task.done)
	task.ctx, task.cancel = context.WithCancel(task.ctx)
	defer task.cancel()

	for {
		select {
		case <-task.ctx.Done():
			log.Printf("Task %d: '%s' Shutting down\n", task.id, task.name)
			return
		case <-time.After(time.Second):
			log.Printf("Task %d: '%s' Running task\n", task.id, task.name)
		}
	}
}

// Add a new task to the group
func (g *TaskGroup) addTask(taskFn TaskFunc, name string, priority int) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second) // Example timeout
	g.addChan <- &Task{
		id:       len(g.tasks),
		name:     name,
		priority: priority,
		ctx:      ctx,
		done:     make(chan struct{}),
	}
}

// Remove a task by index
func (g *TaskGroup) removeTask(index int) {
	g.removeChan <- index
}

// Stop all tasks
func (g *TaskGroup) stop() {
	close(g.shutdownChan)

	// Wait for all tasks to complete
	for _, task := range g.tasks {
		<-task.done
	}
}

// Example scheduled tasks
func task1(ctx context.Context, name string) {
	runTask(ctx, name)
}

func task2(ctx context.Context, name string) {
	runTask(ctx, name)
}

func runTask(ctx context.Context, name string) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("Task %s: Shutting down\n", name)
			return
		case <-ticker.C:
			log.Printf("Task %s: Running task\n", name)
		}
	}
}

func main() {
	group := NewTaskGroup()
	group.start()

	// Adding initial tasks
	group.addTask(task1, "Task 1", 1)
	group.addTask(task2, "Task 2", 2)

	// Let tasks run for a few seconds
	time.Sleep(5 * time.Second)

	// Dynamically add a new task
	group.addTask(task1, "Task 3", 3)

	// Simulate graceful shutdown with priority handling
	log.Println("Initiating graceful shutdown...")
	group.stop()

	log.Println("All tasks have been shut down gracefully.")
}
