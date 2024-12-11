package main

import (
	"container/heap"
	"fmt"
	"sync"
	"time"
)

// Task represents a task with a priority and a function to execute
type Task struct {
	priority int
	exec     func()
	index    int // The index of the task in the heap
}

// A PriorityQueue implements heap.Interface and holds Tasks
type PriorityQueue []*Task

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want higher priority tasks first
	return pq[i].priority > pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	task := x.(*Task)
	task.index = n
	*pq = append(*pq, task)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	task := old[n-1]
	old[n-1] = nil  // prevent memory leak
	task.index = -1 // for safety
	*pq = old[0 : n-1]
	return task
}

// Task Scheduler
type Scheduler struct {
	queue PriorityQueue
	wg    sync.WaitGroup
	quit  chan struct{}
}

// NewScheduler returns a new Scheduler
func NewScheduler() *Scheduler {
	return &Scheduler{
		queue: make(PriorityQueue, 0),
		quit:  make(chan struct{}),
	}
}

// Start initializes the scheduler to process tasks
func (s *Scheduler) Start() {
	heap.Init(&s.queue)

	go func() {
		for {
			select {
			case <-s.quit:
				return
			default:
				s.processTasks()
			}
		}
	}()
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	close(s.quit)
	s.wg.Wait()
}

// AddTask adds a new task to the scheduler with a given priority
func (s *Scheduler) AddTask(priority int, exec func()) {
	s.wg.Add(1)
	heap.Push(&s.queue, &Task{priority: priority, exec: exec})
}

// processTasks executes tasks based on their priority
func (s *Scheduler) processTasks() {
	if s.queue.Len() == 0 {
		time.Sleep(100 * time.Millisecond) // prevent busy waiting
		return
	}
	task := heap.Pop(&s.queue).(*Task)
	go func(t *Task) {
		defer s.wg.Done()
		t.exec()
	}(task)
}

func main() {
	scheduler := NewScheduler()
	scheduler.Start()

	// Add tasks with varying priorities
	scheduler.AddTask(1, func() { fmt.Println("Task with priority 1 executed") })
	scheduler.AddTask(5, func() { fmt.Println("Task with priority 5 executed") })
	scheduler.AddTask(3, func() { fmt.Println("Task with priority 3 executed") })

	// Give some time for tasks to complete
	time.Sleep(2 * time.Second)
	scheduler.Stop()
}
