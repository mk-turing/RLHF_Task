package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Counter struct {
	value int
	mu    sync.Mutex
}

func (c *Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

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

	wg.Wait()
	fmt.Println("Final counter value:", counter.Get())

	// Start workers with error logging and timeout mechanism
	const numWorkers = 5
	errors := make(chan error, numWorkers)
	done := make(chan struct{})

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i, errors, &wg)
	}

	// Handle logged errors and manage completion
	go func() {
		for err := range errors {
			if err != nil {
				fmt.Println("Error occurred:", err)
			}
		}
	}()

	// Wait for workers or timeout
	timeout := time.After(2 * time.Second)
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("All workers completed successfully.")
	case <-timeout:
		fmt.Println("Timeout occurred! Terminating long-running workers.")
		return
	}
}

// worker simulates a task that may fail or exceed a time limit
func worker(id int, errors chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	// Setting a timeout for this worker's operations
	timeout := time.After(1 * time.Second)
	completion := make(chan struct{})

	go func() {
		time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
		// Simulate an error randomly
		if rand.Intn(2) == 0 {
			errors <- fmt.Errorf("worker %d failed", id)
		} else {
			// Indicate completion of worker
			completion <- struct{}{}
		}
	}()

	select {
	case <-completion:
		return
	case <-timeout:
		errors <- fmt.Errorf("worker %d timed out", id)
	}
}
