package main

import (
	"fmt"
	"sync"
	"time"
)

// PriorityMutex provides a priority-based locking mechanism with dynamic priority updates.
type PriorityMutex struct {
	mu              sync.Mutex
	condition       *sync.Cond
	queue           map[chan struct{}]int // Channel -> Priority
	highestPriority int
}

// NewPriorityMutex initializes a new PriorityMutex.
func NewPriorityMutex() *PriorityMutex {
	pm := &PriorityMutex{}
	pm.condition = sync.NewCond(&pm.mu)
	pm.queue = make(map[chan struct{}]int)
	return pm
}

// Lock acquires the lock with the specified priority.
func (pm *PriorityMutex) Lock(priority int) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Add a channel to the queue for this priority level
	ch := make(chan struct{})
	pm.queue[ch] = priority

	// Move the current channel to ensure this priority is at the front if it's higher
	if priority > pm.highestPriority {
		pm.highestPriority = priority
		var newQueue []chan struct{}
		for c, _ := range pm.queue {
			if c == ch {
				newQueue = append(newQueue, c)
			}
		}
		for c, p := range pm.queue {
			if c != ch && p >= priority {
				newQueue = append(newQueue, c)
			}
		}
		for c, p := range pm.queue {
			if c != ch && p < priority {
				newQueue = append(newQueue, c)
			}
		}
		pm.queue = make(map[chan struct{}]int)
		for _, c := range newQueue {
			pm.queue[c] = pm.queue[c]
		}
	}

	// Wait for the signal to acquire the lock
	pm.condition.Wait()
	<-ch
	delete(pm.queue, ch)
}

// Unlock releases the lock, allowing the next goroutine to acquire it.
func (pm *PriorityMutex) Unlock() {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if len(pm.queue) > 0 {
		// Notify the highest priority goroutine
		for ch := range pm.queue {
			close(ch)
			delete(pm.queue, ch)
			break
		}
	}

	pm.condition.Signal()
}

// UpdatePriority updates the priority of a waiting goroutine.
func (pm *PriorityMutex) UpdatePriority(newPriority int) bool {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	for ch, _ := range pm.queue {
		select {
		case <-ch:
			// Channel has already been signaled, cannot update
			return false
		default:
			delete(pm.queue, ch)
			pm.queue[ch] = newPriority

			// Reorder the queue
			if newPriority > pm.highestPriority {
				pm.highestPriority = newPriority
			}

			var newQueue []chan struct{}
			for c, _ := range pm.queue {
				if c == ch {
					newQueue = append(newQueue, c)
				}
			}
			for c, p := range pm.queue {
				if c != ch && p >= newPriority {
					newQueue = append(newQueue, c)
				}
			}
			for c, p := range pm.queue {
				if c != ch && p < newPriority {
					newQueue = append(newQueue, c)
				}
			}
			pm.queue = make(map[chan struct{}]int)
			for _, c := range newQueue {
				pm.queue[c] = pm.queue[c]
			}

			return true
		}
	}
	return false // Goroutine not found in the queue
}

// Example usage of PriorityMutex with dynamic priority updates
func main() {
	pm := NewPriorityMutex()
	var wg sync.WaitGroup

	// Function to simulate work with the PriorityMutex
	work := func(id int, initialPriority int) {
		defer wg.Done()
		pm.Lock(initialPriority)
		defer pm.Unlock()

		fmt.Printf("Goroutine %d with priority %d has acquired the lock\n", id, initialPriority)
		time.Sleep(100 * time.Millisecond) // Simulate work

		// Dynamically update priority
		newPriority := (initialPriority+1)%4 + 1
		pm.UpdatePriority(newPriority)

		fmt.Printf("Goroutine %d updated priority to %d and is releasing the lock\n", id, newPriority)
	}

	// Start multiple goroutines to test the PriorityMutex
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go work(i, i%3+1) // Initial priorities: 1, 2, 3, 1, 2
	}

	wg.Wait()
	fmt.Println("All goroutines have completed.")
}
