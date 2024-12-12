package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Task function simulates a workload.
func task(ctx context.Context, id int) {
	select {
	case <-time.After(time.Duration(rand.Intn(500)) * time.Millisecond):
		fmt.Printf("Task %d completed\n", id)
	case <-ctx.Done():
		fmt.Printf("Task %d cancelled\n", id)
	}
}

// Worker function processes tasks, implementing throttling.
func worker(ctx context.Context, wg *sync.WaitGroup, id int, limiter <-chan struct{}) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case <-limiter: // This channel acts as a semaphore
			task(ctx, id)
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	const maxWorkers = 5
	const maxConcurrentTasks = 3 // Throttling limit
	const totalTasks = 20

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	limiter := make(chan struct{}, maxConcurrentTasks) // Channel for throttling

	// Starting worker goroutines
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go worker(ctx, &wg, i, limiter)
	}

	// Sending tasks
	for i := 0; i < totalTasks; i++ {
		select {
		case limiter <- struct{}{}: // Acquire a spot
			fmt.Printf("Task %d is queued\n", i)
		default:
			fmt.Printf("Task %d is dropped due to load shedding\n", i)
			// If no spot is available, we drop the task
		}

		// Adding a small delay to simulate incoming tasks
		time.Sleep(100 * time.Millisecond)
	}

	// Wait for all workers to finish
	wg.Wait()
}
