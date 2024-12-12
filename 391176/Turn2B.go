package main

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
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
func worker(ctx context.Context, wg *sync.WaitGroup, id int, limiter chan struct{}) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case <-limiter: // This channel acts as a semaphore
			task(ctx, id)
			limiter <- struct{}{}
		}
	}
}

// AdjustRateLimit dynamically adjusts the rate limit based on CPU usage.
func adjustRateLimit(ctx context.Context, limiter *chan struct{}, maxConcurrentTasks *int32) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Calculate current CPU usage
			var gomaxprocs = runtime.GOMAXPROCS(0)
			var avgLoad runtime.LoadAvg
			runtime.ReadLoadAvg(&avgLoad)
			currentLoad := avgLoad.Load1 / float64(gomaxprocs)

			// Dynamic rate limiting logic
			var newMax int32
			if currentLoad > 1.5 {
				// High load: reduce concurrency
				newMax = max(1, int32(float64(*maxConcurrentTasks)*0.7))
			} else if currentLoad > 0.5 {
				// Medium load: maintain concurrency
				newMax = *maxConcurrentTasks
			} else {
				// Low load: increase concurrency
				newMax = min(10, int32(float64(*maxConcurrentTasks)*1.3))
			}

			// Resize the limiter channel if necessary
			if newMax != *maxConcurrentTasks {
				oldCap := cap(*limiter)
				newLimiter := make(chan struct{}, newMax)
				for i := 0; i < oldCap; i++ {
					select {
					case <-*limiter:
						newLimiter <- struct{}{}
					default:
					}
				}
				*limiter = newLimiter
				atomic.StoreInt32(maxConcurrentTasks, newMax)
				fmt.Printf("Adjusted rate limit to %d tasks\n", newMax)
			}
		}
	}
}

func max(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func min(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func main() {
	rand.Seed(time.Now().UnixNano())

	const maxWorkers = 5
	var maxConcurrentTasks int32 = 5 // Initial throttling limit
	const totalTasks = 100

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	limiter := make(chan struct{}, maxConcurrentTasks) // Channel for throttling

	// Start adjusting rate limit in a goroutine
	go adjustRateLimit(ctx, &limiter, &maxConcurrentTasks)

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
