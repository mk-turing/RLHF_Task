package main

import (
	"container/heap"
	"fmt"
	"sync"
	"time"
)

// PriorityMutex provides a high-precision locking mechanism with priority support.
type PriorityMutex struct {
	mu     sync.Mutex
	cond   *sync.Cond
	queue  priorityQueue // Min-heap based on goroutine priority
	locked bool
}

// A GoroutineRequest represents a request for the mutex with an associated priority.
type GoroutineRequest struct {
	id       int // Goroutine ID
	priority int // Priority level
	index    int // Index in the priority queue
}

// A priorityQueue implements heap.Interface and holds requests.
type priorityQueue []*GoroutineRequest

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority // Min-heap based on priority
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	request := x.(*GoroutineRequest)
	request.index = n
	*pq = append(*pq, request)
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	request := old[n-1]
	old[n-1] = nil     // avoid memory leak
	request.index = -1 // for safety
	*pq = old[0 : n-1]
	return request
}

// NewPriorityMutex initializes a new PriorityMutex.
func NewPriorityMutex() *PriorityMutex {
	pm := &PriorityMutex{}
	pm.cond = sync.NewCond(&pm.mu)
	return pm
}

// Lock locks the PriorityMutex, blocking if necessary until it can be acquired.
func (pm *PriorityMutex) Lock(id int, priority int) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	request := &GoroutineRequest{id: id, priority: priority}
	heap.Push(&pm.queue, request)

	for pm.locked || (pm.queue[0] != request) {
		pm.cond.Wait()
	}

	pm.locked = true
	heap.Pop(&pm.queue) // Remove the request from the queue
}

// Unlock unlocks the PriorityMutex, allowing others to acquire it.
func (pm *PriorityMutex) Unlock() {
	pm.mu.Lock()
	pm.locked = false
	pm.cond.Broadcast() // Wake up waiting goroutines
	pm.mu.Unlock()
}

// Example usage of PriorityMutex
func main() {
	pm := NewPriorityMutex()
	var wg sync.WaitGroup

	// Function to simulate work with the PriorityMutex
	work := func(id int, priority int) {
		defer wg.Done()
		pm.Lock(id, priority)
		defer pm.Unlock()

		fmt.Printf("Goroutine %d with priority %d has acquired the lock\n", id, priority)
		time.Sleep(100 * time.Millisecond) // Simulate work
		fmt.Printf("Goroutine %d is releasing the lock\n", id)
	}

	// Start multiple goroutines with varying priorities
	goroutines := []struct {
		id       int
		priority int
	}{
		{1, 1}, // Low priority
		{2, 3}, // High priority
		{3, 2}, // Medium priority
		{4, 1}, // Low priority
		{5, 5}, // Lowest priority
	}

	for _, g := range goroutines {
		wg.Add(1)
		go work(g.id, g.priority)
	}

	wg.Wait()
	fmt.Println("All goroutines have completed.")
}
