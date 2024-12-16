package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Counter represents a shared counter that can be incremented concurrently
type Counter struct {
	value int
	mu    sync.Mutex
}

// Inc increments the counter by 1
func (c *Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

// Get returns the current value of the counter
func (c *Counter) Get() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func main() {
	var wg sync.WaitGroup
	counter := &Counter{}
	numGoroutines := 100
	numIncrements := 1000

	// Start goroutines to increment the counter
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numIncrements; j++ {
				counter.Inc()
			}
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Print the final counter value
	fmt.Println("Final counter value:", counter.Get())

	// Simulate an error handling in concurrent scenario
	err := simulateConcurrentError()
	if err != nil {
		fmt.Println("Error occurred:", err)
	} else {
		fmt.Println("Operation completed successfully.")
	}
}

// simulateConcurrentError demonstrates error handling in a concurrent scenario
func simulateConcurrentError() error {
	const numWorkers = 5
	var wg sync.WaitGroup
	errors := make(chan error, numWorkers)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			// Simulate an error
			if rand.Intn(2) == 0 {
				errors <- fmt.Errorf("worker %d failed", i)
			}
		}()
	}

	wg.Wait()

	close(errors)
	for err := range errors {
		return err
	}
	return nil
}
