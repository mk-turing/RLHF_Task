package main

import (
	"fmt"
	"sync"
	"time"
)

// PrecisionMutex provides a high-precision locking mechanism.
type PrecisionMutex struct {
	mu        sync.Mutex
	condition *sync.Cond
	locked    bool
}

// NewPrecisionMutex initializes a new PrecisionMutex.
func NewPrecisionMutex() *PrecisionMutex {
	pm := &PrecisionMutex{}
	pm.condition = sync.NewCond(&pm.mu)
	return pm
}

// Lock locks the PrecisionMutex, blocking if necessary until it can be acquired.
func (pm *PrecisionMutex) Lock() {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	for pm.locked {
		pm.condition.Wait()
	}
	pm.locked = true
}

// Unlock unlocks the PrecisionMutex, allowing others to acquire it.
func (pm *PrecisionMutex) Unlock() {
	pm.mu.Lock()
	pm.locked = false
	pm.condition.Broadcast() // Wake up waiting goroutines
	pm.mu.Unlock()
}

// IsLocked returns true if the mutex is currently locked.
func (pm *PrecisionMutex) IsLocked() bool {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	return pm.locked
}

// Example usage of PrecisionMutex
func main() {
	pm := NewPrecisionMutex()
	var wg sync.WaitGroup

	// Function to simulate work with the PrecisionMutex
	work := func(id int) {
		defer wg.Done()
		pm.Lock()
		defer pm.Unlock()

		fmt.Printf("Goroutine %d has acquired the lock\n", id)
		time.Sleep(100 * time.Millisecond) // Simulate work
		fmt.Printf("Goroutine %d is releasing the lock\n", id)
	}

	// Start multiple goroutines to test the PrecisionMutex
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go work(i)
	}

	wg.Wait()
	fmt.Println("All goroutines have completed.")
}
