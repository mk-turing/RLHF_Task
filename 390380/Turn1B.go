package main

import (
	"fmt"
	"sync"
)

// ReadWriteMutex allows multiple read goroutines to run concurrently
// but only one write goroutine at a time.
type ReadWriteMutex struct {
	mu      sync.Mutex
	readers int
	writing bool
	cond    *sync.Cond
}

// NewReadWriteMutex creates a new ReadWriteMutex.
func NewReadWriteMutex() *ReadWriteMutex {
	return &ReadWriteMutex{
		cond: sync.NewCond(&sync.Mutex{}),
	}
}

// Lock acquires the write lock.
func (rwm *ReadWriteMutex) Lock() {
	rwm.mu.Lock()
	defer rwm.mu.Unlock()

	for rwm.writing {
		rwm.cond.Wait()
	}
	rwm.writing = true
}

// Unlock releases the write lock.
func (rwm *ReadWriteMutex) Unlock() {
	rwm.mu.Lock()
	defer rwm.mu.Unlock()

	rwm.writing = false
	rwm.cond.Broadcast()
}

// RLock acquires the read lock.
func (rwm *ReadWriteMutex) RLock() {
	rwm.mu.Lock()
	defer rwm.mu.Unlock()

	for rwm.writing {
		rwm.cond.Wait()
	}
	rwm.readers++
}

// RUnlock releases the read lock.
func (rwm *ReadWriteMutex) RUnlock() {
	rwm.mu.Lock()
	defer rwm.mu.Unlock()

	rwm.readers--
	if rwm.readers == 0 {
		rwm.cond.Broadcast()
	}
}

func main() {
	var rwm ReadWriteMutex
	var sharedData int

	// Worker function to simulate read operations
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rwm.RLock()
			defer rwm.RUnlock()
			fmt.Printf("Reading shared data: %d\n", sharedData)
		}()
	}

	// Worker function to simulate write operations
	wg.Add(1)
	go func() {
		defer wg.Done()
		rwm.Lock()
		defer rwm.Unlock()
		sharedData = 42
		fmt.Printf("Updated shared data to: %d\n", sharedData)
	}()

	wg.Wait()
}
