package main

import (
	"context"
	"fmt"
	"log"
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

	// Simulate concurrent tasks with potential errors and timeouts
	startTime := time.Now()
	err := simulateConcurrentTasksWithTimeout(10 * time.Second) // Timeout after 10 seconds
	duration := time.Since(startTime)
	log.Printf("Operation completed in %s. Error: %v\n", duration, err)
}

// simulateConcurrentTasksWithTimeout demonstrates error logging and timeout handling
func simulateConcurrentTasksWithTimeout(timeout time.Duration) error {
	const numWorkers = 5
	var wg sync.WaitGroup
	errors := make(chan error, numWorkers)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				log.Printf("Worker %d timed out\n", id)
				return
			default:
				time.Sleep(time.Duration(500) * time.Millisecond)
				// Simulate an error
				if 1 == 0 {
					errors <- fmt.Errorf("worker %d failed", id)
					log.Printf("Worker %d encountered an error: %v\n", id)
				}
			}
		}(i)
	}

	wg.Wait()

	close(errors)
	for err := range errors {
		return err
	}
	return nil
}
