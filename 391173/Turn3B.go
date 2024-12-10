package main

import (
	"context"
	"log"
	"sync"
	"time"
)

// Define a type for task functions
type taskFunc func(context.Context)

// Manage a group of tasks
type taskGroup struct {
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	tasks   []taskFunc
	name    string
	mu      sync.Mutex
	channel chan struct{} // Optional: for additional control
}

func (g *taskGroup) start() {
	for _, task := range g.tasks {
		g.wg.Add(1)
		go func(t taskFunc) {
			defer g.wg.Done()
			t(g.ctx)
		}(task)
	}
}

func (g *taskGroup) addTask(task taskFunc) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.tasks = append(g.tasks, task)
	g.wg.Add(1)
	go func(t taskFunc) {
		defer g.wg.Done()
		t(g.ctx)
	}(task)
}

func (g *taskGroup) removeTask(index int) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if index >= 0 && index < len(g.tasks) {
		g.tasks = append(g.tasks[:index], g.tasks[index+1:]...)
		g.wg.Done()
	}
}

func (g *taskGroup) wait() {
	g.wg.Wait()
	close(g.channel) // Optional: signal all tasks through the channel
	log.Printf("Task group '%s' finished.\n", g.name)
}

func (g *taskGroup) stop() {
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

func task3(ctx context.Context) {
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
}

func main() {
	// Create task groups
	group1 := &taskGroup{
		tasks:   []taskFunc{task1},
		name:    "Group 1",
		channel: make(chan struct{}),
	}

	group2 := &taskGroup{
		tasks: []taskFunc{task2, task3},
		name:  "Group 2",
	}

	// Start task groups
	ctx, cancelAll := context.WithCancel(context.Background())
	defer cancelAll()

	group1.ctx, group1.cancel = context.WithCancel(ctx)
	group2.ctx, group2.cancel = context.WithCancel(ctx)

	group1.start()
	group2.start()

	// Simulate runtime task additions and removals
	time.Sleep(2 * time.Second)
	log.Println("Adding Task 4 to Group 1")
	group1.addTask(func(ctx context.Context) {
		ticker := time.NewTicker(4 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				log.Println("Task 4: Shutting down")
				return
			case <-ticker.C:
				log.Println("Task 4: Running task")
			}
		}
	})

	time.Sleep(3 * time.Second)
	log.Println("Removing Task 2 from Group 2")
	group2.removeTask(0)

	// Simulate graceful shutdown after some time
	time.Sleep(5 * time.Second)
	log.Println("Initiating graceful shutdown...")

	// Stop task groups
	group1.stop()
	group2.stop()

	log.Println("All tasks have been shut down gracefully.")
}
