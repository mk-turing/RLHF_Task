package main

import (
	"fmt"
	"sync"
	"time"
)

// PriorityMutex provides a priority-based locking mechanism.
type PriorityMutex struct {
	mu              sync.Mutex
	condition       *sync.Cond
	queue           []chan struct{}
	highestPriority int
}

// NewPriorityMutex initializes a new PriorityMutex.
func NewPriorityMutex() *PriorityMutex {
	pm := &PriorityMutex{}
	pm.condition = sync.NewCond(&pm.mu)
	return pm
}

// Lock acquires the lock with the specified priority.
func (pm *PriorityMutex) Lock(priority int) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Add a channel to the queue for this priority level
	ch := make(chan struct{})
	pm.queue = append(pm.queue, ch)

	// Move the current queue to ensure this priority is at the front if it's higher
	if priority > pm.highestPriority {
		pm.highestPriority = priority
		for i := range pm.queue {
			if pm.queue[i] == ch {
				pm.queue = append(pm.queue[:i], pm.queue[i+1:]...)
				pm.queue = append([]chan struct{}{ch}, pm.queue...)
				break
			}
		}
	}

	// Wait for the signal to acquire the lock
	pm.condition.Wait()
	<-ch
}

// Unlock releases the lock, allowing the next goroutine to acquire it.
func (pm *PriorityMutex) Unlock() {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if len(pm.queue) > 0 {
		// Notify the highest priority goroutine
		ch := pm.queue[0]
		pm.queue = pm.queue[1:]
		close(ch)
	}

	pm.condition.Signal()
}

// Example usage of PriorityMutex
func main() {
	pm := NewPriorityMutex()
	var wg sync.WaitGroup

	// Function to simulate work with the PriorityMutex
	work := func(id int, priority int) {
		defer wg.Done()
		pm.Lock(priority)
		defer pm.Unlock()

		fmt.Printf("Goroutine %d with priority %d has acquired the lock\n", id, priority)
		time.Sleep(100 * time.Millisecond) // Simulate work
		fmt.Printf("Goroutine %d with priority %d is releasing the lock\n", id, priority)
	}

	// Start multiple goroutines to test the PriorityMutex
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go work(i, i%3+1) // Priorities: 1, 2, 3, 1, 2
	}

	wg.Wait()
	fmt.Println("All goroutines have completed.")
}
