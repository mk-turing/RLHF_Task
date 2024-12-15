package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// PriorityEntry represents a goroutine waiting for the lock with its associated priority and deadline.
type PriorityEntry struct {
	priority   int
	ctx        context.Context
	ch         chan struct{}
	deadline   time.Time
	cancelFunc context.CancelFunc
}

// DynamicPriorityMutex provides a dynamic priority-based locking mechanism with deadline support.
type DynamicPriorityMutex struct {
	mu        sync.Mutex
	condition *sync.Cond
	queue     []PriorityEntry
	ticker    *time.Ticker
}

// NewDynamicPriorityMutex initializes a new DynamicPriorityMutex.
func NewDynamicPriorityMutex() *DynamicPriorityMutex {
	pm := &DynamicPriorityMutex{}
	pm.condition = sync.NewCond(&pm.mu)
	pm.ticker = time.NewTicker(time.Millisecond * 100) // Check deadlines every 100ms
	go pm.checkDeadlines()
	return pm
}

// Lock acquires the lock with the specified priority and context.
func (pm *DynamicPriorityMutex) Lock(ctx context.Context) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	priority := rand.Intn(5) + 1 // Random priority between 1 and 5
	deadline, _ := ctx.Deadline()
	if deadline.IsZero() {
		// Set a default deadline if none is provided
		deadline = time.Now().Add(time.Duration(rand.Intn(1000)) * time.Millisecond)
	}

	// Create a new entry and add it to the queue
	ch := make(chan struct{})
	entry := PriorityEntry{priority: priority, ctx: ctx, ch: ch, deadline: deadline}
	pm.queue = append(pm.queue, entry)

	// Sort based on priority (higher priority first)
	//pm.sortQueue()

	// Wait for the signal to acquire the lock or until context is canceled
	for {
		select {
		case <-entry.ctx.Done():
			fmt.Printf("Goroutine %p canceled\n", ctx)
			return
		//case <-pm.condition.C:
		//	if pm.queue[0].ch == ch {
		//		break
		//	}
		default:
			if entry.deadline.Before(time.Now()) {
				fmt.Printf("Goroutine %p deadline exceeded\n", ctx)
				return
			}
		}
	}

	// Remove entry from the queue and proceed
	pm.queue = pm.queue[1:]
	close(ch)
}

// Unlock releases the lock, allowing the next goroutine to acquire it.
func (pm *DynamicPriorityMutex) Unlock() {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if len(pm.queue) > 0 {
		ch := pm.queue[0].ch
		pm.queue = pm.queue[1:]
		close(ch)
	}

	pm.condition.Signal()
}

// checkDeadlines periodically checks if any entries have exceeded their deadlines.
func (pm *DynamicPriorityMutex) checkDeadlines() {
	for range pm.ticker.C {
		pm.mu.Lock()
		defer pm.mu.Unlock()

		now := time.Now()
		i := 0
		for i < len(pm.queue) {
			if pm.queue[i].deadline.Before(now) {
				fmt.Printf("Deadline for goroutine %p exceeded\n", pm.queue[i].ctx)
				pm.queue[i].cancelFunc()
				pm.queue = append(pm.queue[:i], pm.queue[i+1:]...)
			} else {
				i++
			}
		}

		pm.condition.Broadcast() // Reschedule queue check
	}
}

// Example usage of DynamicPriorityMutex
func main() {
	pm := NewDynamicPriorityMutex()
	var wg sync.WaitGroup

	// Function to simulate work with the DynamicPriorityMutex
	work := func(id int) {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(rand.Intn(500)+100)*time.Millisecond)
		defer cancel()

		pm.Lock(ctx)
		defer pm.Unlock()

		fmt.Printf("Goroutine %d with priority %d has acquired the lock\n", id, rand.Intn(5)+1)

		// Simulate work
		time.Sleep(100 * time.Millisecond)

		fmt.Printf("Goroutine %d is releasing the lock\n", id)
	}

	// Start multiple goroutines to test the DynamicPriorityMutex
	numGoroutines := 10
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go work(i)
	}

	wg.Wait()
	fmt.Println("All goroutines have completed.")
}
