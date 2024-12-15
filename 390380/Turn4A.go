package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// PriorityEntry represents a goroutine waiting for the lock with its associated priority and deadline.
type PriorityEntry struct {
	priority int
	deadline time.Time
	ch       chan struct{}
}

// DeadlinePriorityMutex provides a locking mechanism based on priority and deadlines.
type DeadlinePriorityMutex struct {
	mu        sync.Mutex
	condition *sync.Cond
	queue     []PriorityEntry
}

// NewDeadlinePriorityMutex initializes a new DeadlinePriorityMutex.
func NewDeadlinePriorityMutex() *DeadlinePriorityMutex {
	dpm := &DeadlinePriorityMutex{}
	dpm.condition = sync.NewCond(&dpm.mu)
	return dpm
}

// Lock tries to acquire the lock with a specified priority and deadline.
func (dpm *DeadlinePriorityMutex) Lock(priority int, deadline time.Time) bool {
	dpm.mu.Lock()
	defer dpm.mu.Unlock()

	ch := make(chan struct{})
	entry := PriorityEntry{priority: priority, deadline: deadline, ch: ch}
	dpm.queue = append(dpm.queue, entry)

	// Wait for the signal to acquire the lock
	for {
		// Remove expired entries from the front of the queue
		dpm.cleanQueue()

		// Sort based on priority (higher priority first)
		dpm.sortQueue()

		if len(dpm.queue) > 0 && dpm.queue[0].ch == ch {
			// Acquired lock
			break
		}

		if time.Now().After(deadline) {
			// Deadline exceeded; clean up and exit
			close(ch)
			return false
		}

		dpm.condition.Wait()
	}

	// Remove entry from the queue and proceed
	dpm.queue = dpm.queue[1:]
	close(ch)
	return true
}

// Unlock releases the lock, allowing the next goroutine to acquire it.
func (dpm *DeadlinePriorityMutex) Unlock() {
	dpm.mu.Lock()
	defer dpm.mu.Unlock()
	dpm.condition.Signal()
}

// cleanQueue removes expired entries from the queue.
func (dpm *DeadlinePriorityMutex) cleanQueue() {
	currentTime := time.Now()
	for len(dpm.queue) > 0 && dpm.queue[0].deadline.Before(currentTime) {
		dpm.queue = dpm.queue[1:] // Remove expired entries
	}
}

// sortQueue sorts the queue by priority (highest first).
func (dpm *DeadlinePriorityMutex) sortQueue() {
	n := len(dpm.queue)
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			if dpm.queue[i].priority < dpm.queue[j].priority {
				dpm.queue[i], dpm.queue[j] = dpm.queue[j], dpm.queue[i]
			}
		}
	}
}

// Example usage of DeadlinePriorityMutex
func main() {
	dpm := NewDeadlinePriorityMutex()
	var wg sync.WaitGroup

	// Function to simulate work with the DeadlinePriorityMutex
	work := func(id int, priority int, deadline time.Time) {
		defer wg.Done()
		acquired := dpm.Lock(priority, deadline)
		if !acquired {
			fmt.Printf("Goroutine %d with priority %d could not acquire the lock before the deadline\n", id, priority)
			return
		}

		fmt.Printf("Goroutine %d with priority %d has acquired the lock\n", id, priority)

		// Simulate work
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

		fmt.Printf("Goroutine %d with priority %d is releasing the lock\n", id, priority)
		dpm.Unlock()
	}

	// Start multiple goroutines to test the DeadlinePriorityMutex
	numGoroutines := 8
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		// Generate a random deadline (current time + random duration up to 300 ms)
		deadline := time.Now().Add(time.Duration(rand.Intn(300)) * time.Millisecond)
		go work(i, rand.Intn(5)+1, deadline) // Initial random priorities between 1 and 5
	}

	wg.Wait()
	fmt.Println("All goroutines have completed.")
}
