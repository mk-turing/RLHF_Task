package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type RateLimiter struct {
	limit       int64         // Max number of requests per window
	interval    time.Duration // Length of the window
	windowStart time.Time
	remaining   int64
	mutex       sync.Mutex
	counter     *atomic.Int64
}

// NewRateLimiter creates a new rate limiter with specified limits and intervals.
func NewRateLimiter(limit int64, interval time.Duration) *RateLimiter {
	counter := new(atomic.Int64) // Fix: Initialize the atomic counter
	return &RateLimiter{
		limit:       limit,
		interval:    interval,
		windowStart: time.Now(),
		remaining:   limit,
		counter:     counter,
	}
}

// Acquire attempts to acquire a request permit from the rate limiter.
// It waits for the next window if the limit has been reached in the current window.
func (rl *RateLimiter) Acquire(ctx context.Context) (context.Context, context.CancelFunc, error) {
	if ctx == nil {
		return nil, nil, fmt.Errorf("context is nil")
	}

	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	// Check if it's the start of a new window
	if time.Now().Sub(rl.windowStart) >= rl.interval {
		rl.windowStart = time.Now()
		rl.remaining = rl.limit
	}

	// Block until a permit is available
	for {
		select {
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		default:
		}
		if rl.remaining > 0 {
			rl.remaining--
			rl.counter.Add(1)
			break
		}
		// Wait until the next window
		delta := rl.windowStart.Add(rl.interval).Sub(time.Now())
		if delta > 0 {
			time.Sleep(delta)
		}
	}

	// Return a new context with a cancel function
	return ctx, nil, nil
}

// Call with rate limit, retries, and context timeout.
func callWithRateLimit(ctx context.Context, url string, retries int) error {
	rl := NewRateLimiter(3, 5*time.Second)

	err := error(nil)
	for retry := 0; retry <= retries; retry++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		newCtx, _, err := rl.Acquire(ctx)
		if err != nil {
			fmt.Printf("Failed to acquire limit: %v\n", err)
			return err
		}

		err = makeApiCall(newCtx, url)
		if err == nil {
			fmt.Printf("Call to %s succeeded after %d retries\n", url, retry)
			return nil
		}

		select {
		case <-newCtx.Done():
			fmt.Printf("Call to %s cancelled: %v\n", url, newCtx.Err())
			return newCtx.Err()
		default:
		}

		fmt.Printf("Retried: err: %v\n", err)
		// Add retry strategy as needed
		// time.Sleep(1*time.Second)
	}

	fmt.Printf("Call to %s failed after %d retries: %v\n", url, retries, err)
	return err
}

// makeApiCall makes an HTTP request with a given context and simulates a delay.
func makeApiCall(ctx context.Context, url string) error {
	// Simulate a delay before making the request
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		time.Sleep(1 * time.Second)
	}

	// Randomly fail with a 30% probability for demonstration
	if float64(rand.Intn(10)) < 3 {
		return fmt.Errorf("simulated failure for %s", url)
	}

	return nil
}

func main() {
	urls := []string{
		"https://example.com/api1",
		"https://example.com/api2",
		"https://example.com/api3",
	}
	maxRetries := 2
	var wg sync.WaitGroup

	// Start the API calls with rate limiting, retries, and context timeouts
	for _, url := range urls {
		// Initialize context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		// Defer cancel inside the goroutine, after starting it
		wg.Add(1)
		go func(url string, maxRetries int, ctx context.Context, cancel context.CancelFunc) {
			defer wg.Done()
			defer cancel() // Ensure cancel is deferred properly inside the goroutine

			err := callWithRateLimit(ctx, url, maxRetries)
			if err != nil {
				fmt.Printf("Error calling %s: %v\n", url, err)
			}
		}(url, maxRetries, ctx, cancel) // Pass cancel as a parameter to the goroutine
	}

	wg.Wait() // Wait for all goroutines to complete
}
