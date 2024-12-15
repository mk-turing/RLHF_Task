package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// PriorityEntry represents a goroutine waiting for the lock with its associated priority.
type PriorityEntry struct {
	priority int
	ch       chan struct{}
}

// DynamicPriorityMutex provides a dynamic priority-based locking mechanism.
type DynamicPriorityMutex struct {
	mu        sync.Mutex
	condition *sync.Cond
	queue     []PriorityEntry
}

// NewDynamicPriorityMutex initializes a new DynamicPriorityMutex.
func NewDynamicPriorityMutex() *DynamicPriorityMutex {
	pm := &DynamicPriorityMutex{}
	pm.condition = sync.NewCond(&pm.mu)
	return pm
}

// Lock acquires the lock with the specified priority.
func (pm *DynamicPriorityMutex) Lock(priority int) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Create a new entry and add it to the queue
	ch := make(chan struct{})
	entry := PriorityEntry{priority: priority, ch: ch}
	pm.queue = append(pm.queue, entry)

	// Wait for the signal to acquire the lock
	for {
		// Sort based on priority (higher priority first)
		pm.sortQueue()
		if pm.queue[0].ch == ch {
			break
		}
		pm.condition.Wait()
	}
	// Remove entry from the queue and proceed
	pm.queue = pm.queue[1:]
	close(ch)
}

// Unlock releases the lock, allowing the next goroutine to acquire it.
func (pm *DynamicPriorityMutex) Unlock() {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.condition.Signal()
}

// ChangePriority updates the priority of the calling goroutine while waiting.
func (pm *DynamicPriorityMutex) ChangePriority(oldPriority, newPriority int) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	for i, entry := range pm.queue {
		if entry.priority == oldPriority {
			pm.queue[i].priority = newPriority
			break
		}
	}
}

// sortQueue sorts the queue by priority (highest first).
func (pm *DynamicPriorityMutex) sortQueue() {
	n := len(pm.queue)
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			if pm.queue[i].priority < pm.queue[j].priority {
				pm.queue[i], pm.queue[j] = pm.queue[j], pm.queue[i]
			}
		}
	}
}

// Example usage of DynamicPriorityMutex
func main() {
	pm := NewDynamicPriorityMutex()
	var wg sync.WaitGroup

	// Function to simulate work with the DynamicPriorityMutex
	work := func(id int, priority int) {
		defer wg.Done()
		pm.Lock(priority)
		defer pm.Unlock()

		fmt.Printf("Goroutine %d with priority %d has acquired the lock\n", id, priority)

		// Simulate work and random priority change
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

		// Randomly change priority
		newPriority := rand.Intn(5) + 1 // New priority between 1 and 5
		fmt.Printf("Goroutine %d changing priority from %d to %d\n", id, priority, newPriority)
		pm.ChangePriority(priority, newPriority)

		// Simulate more work
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("Goroutine %d with priority %d is releasing the lock\n", id, newPriority)
	}

	// Start multiple goroutines to test the DynamicPriorityMutex
	numGoroutines := 5
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go work(i, rand.Intn(5)+1) // Initial random priorities between 1 and 5
	}

	wg.Wait()
	fmt.Println("All goroutines have completed.")
}
