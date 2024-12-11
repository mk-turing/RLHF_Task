package _91164

import (
	"sync"
	"testing"
	"time"
)

type Counter struct {
	mu    sync.Mutex
	value int
}

func (c *Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *Counter) Dec() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.value > 0 {
		c.value--
	}
}

func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func StartStatefulGoroutine(counter *Counter) chan int {
	done := make(chan int)

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				counter.Inc()
			}
		}
	}()

	return done
}

func TestStatefulGoroutine(t *testing.T) {
	counter := &Counter{}
	done := StartStatefulGoroutine(counter)

	// Allow the goroutine to run for a moment
	time.Sleep(10 * time.Millisecond)

	// Test current state
	if counter.Value() != 0 {
		t.Errorf("Expected counter to be 0, got %d", counter.Value())
	}

	// Simulate external actions
	for i := 0; i < 5; i++ {
		counter.Inc()
	}

	// Allow the goroutine to update its state
	time.Sleep(10 * time.Millisecond)

	// Test updated state
	if counter.Value() != 5 {
		t.Errorf("Expected counter to be 5, got %d", counter.Value())
	}

	// Stop the goroutine
	close(done)

	// Wait for goroutine to stop
	<-done

	// Final state check
	if counter.Value() != 5 {
		t.Errorf("Expected counter to be 5 after stopping, got %d", counter.Value())
	}
}
