package _91164

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// Counter holds the state and provides methods to manipulate it safely.
type Counter struct {
	mu    sync.Mutex
	value int
}

// Increment increments the counter's value.
func (c *Counter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

// GetValue returns the current value of the counter.
func (c *Counter) GetValue() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

// StressTest simulates high-concurrency scenarios.
func StressTest(t *testing.T, numGoroutines int, incrementsPerGoroutine int) {
	counter := &Counter{}
	var wg sync.WaitGroup

	start := time.Now()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				counter.Increment()
			}
		}()
	}

	wg.Wait()

	duration := time.Since(start)

	fmt.Printf("Total increments: %d, Duration: %s\n", numGoroutines*incrementsPerGoroutine, duration)

	expectedValue := numGoroutines * incrementsPerGoroutine
	if counter.GetValue() != expectedValue {
		t.Errorf("Expected counter value %d, got %d", expectedValue, counter.GetValue())
	}
}

// TestConcurrentIncrement tests the Counter under high concurrency.
func TestConcurrentIncrement(t *testing.T) {
	// Parameters for testing
	numGoroutines := 100
	incrementsPerGoroutine := 1000

	StressTest(t, numGoroutines, incrementsPerGoroutine)
}
