package _91164

import (
	"sync"
	"testing"
	"time"
)

// Stateful Goroutine
type Counter struct {
	mu      sync.Mutex
	value   int
	stopped bool
}

func (c *Counter) Start(ch chan string) {
	go func() {
		for {
			select {
			case msg := <-ch:
				if msg == "inc" {
					c.mu.Lock()
					c.value++
					c.mu.Unlock()
				} else if msg == "stop" {
					c.mu.Lock()
					c.stopped = true
					c.mu.Unlock()
					return
				}
			default:
				time.Sleep(10 * time.Millisecond) // Simulate work
			}
		}
	}()
}

func (c *Counter) GetValue() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func (c *Counter) IsStopped() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.stopped
}

// Test for Counter
func TestCounter(t *testing.T) {
	ct := &Counter{}
	ch := make(chan string)

	ct.Start(ch)

	// Increment the counter
	ch <- "inc"
	time.Sleep(50 * time.Millisecond) // Wait for goroutine to process
	if ct.GetValue() != 1 {
		t.Errorf("Expected counter value to be 1, got %d", ct.GetValue())
	}

	// Increment again
	ch <- "inc"
	time.Sleep(50 * time.Millisecond)
	if ct.GetValue() != 2 {
		t.Errorf("Expected counter value to be 2, got %d", ct.GetValue())
	}

	// Stop the counter
	ch <- "stop"
	time.Sleep(50 * time.Millisecond) // Wait for goroutine to stop
	if !ct.IsStopped() {
		t.Error("Expected counter to be stopped")
	}
}
