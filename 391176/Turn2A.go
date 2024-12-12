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

// Worker function processes tasks, using a semaphore for throttling.
func worker(ctx context.Context, wg *sync.WaitGroup, id int, limiter <-chan struct{}) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case <-limiter: // Acquire a spot in the limiter
			task(ctx, id)
		}
	}
}

// SimulateLoad monitors and adjusts the load dynamically.
func simulateLoad(ctx context.Context, loadCh chan<- int, adjustmentInterval time.Duration) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Simulate a load value between 1 and 10
			loadValue := rand.Intn(10) + 1
			loadCh <- loadValue
			time.Sleep(adjustmentInterval)
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	const maxWorkers = 5
	const totalTasks = 20
	loadCh := make(chan int)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var wg sync.WaitGroup

	// Starting worker goroutines
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		limiter := make(chan struct{}, 3) // Initial concurrency limit

		go worker(ctx, &wg, i, limiter)

		// Adjust the limiter based on system load
		go func() {
			for {
				select {
				case load := <-loadCh:
					// Adjust the concurrency limit based on the load
					newLimit := 3 - (load - 1) // Adjust limit based on load
					if newLimit < 1 {
						newLimit = 1
					} else if newLimit > 3 {
						newLimit = 3
					}
					fmt.Printf("Adjusting concurrency limit to %d based on load %d\n", newLimit, load)

					// Update the semaphore (limiter)
					limiter = make(chan struct{}, newLimit)
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	// Start load simulation
	go simulateLoad(ctx, loadCh, 2*time.Second)

	// Sending tasks
	for i := 0; i < totalTasks; i++ {
		select {
		case <-ctx.Done():
			break
		default:
			// Try to acquire a spot in the limiter
			if len(limiter) < cap(limiter) {
				limiter <- struct{}{} // Acquire a spot
				fmt.Printf("Task %d is queued\n", i)
			} else {
				fmt.Printf("Task %d is dropped due to load shedding\n", i)
			}
		}
		time.Sleep(100 * time.Millisecond) // Simulate incoming tasks
	}

	// Wait for all workers to finish
	wg.Wait()
}
