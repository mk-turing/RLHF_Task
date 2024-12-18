package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin/concurrencylimiter"
	"github.com/gin-gonic/gin/concurrencylimiter/limiter"
)

// retryApiCall attempts to make an HTTP request with a rate limit and retries.
func retryApiCall(ctx context.Context, url string, delay time.Duration, maxRetries int, result chan<- string) {
	err := error(nil)
	for retry := 0; retry <= maxRetries; retry++ {
		// Acquire a rate limit ticket
		err := l.Wait(ctx)
		if err != nil {
			result <- fmt.Sprintf("Call to %s cancelled due to rate limiting: %v", url, err)
			return
		}
		select {
		case <-ctx.Done():
			result <- fmt.Sprintf("Call to %s cancelled: %v", url, ctx.Err())
			return
		default:
		}

		err = makeApiCall(ctx, url)
		if err == nil {
			result <- fmt.Sprintf("Call to %s succeeded after %d retries", url, retry)
			return
		}

		select {
		case <-ctx.Done():
			result <- fmt.Sprintf("Call to %s cancelled: %v", url, ctx.Err())
			return
		default:
		}
		// Backoff strategy: exponential backoff with a minimum delay
		time.Sleep(time.Duration(pow(2.0, retry)) * time.Second)
	}
	result <- fmt.Sprintf("Call to %s failed after %d retries: %v", url, maxRetries, err)
}

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
	// Set up the rate limiter
	l := limiter.NewFixedWindowLimiter(5, time.Second)

	urls := []string{
		"https://example.com/api1",
		"https://example.com/api2",
		"https://example.com/api3",
	}
	delays := []time.Duration{
		3 * time.Second,
		4 * time.Second,
		2 * time.Second,
	}
	maxRetries := 2
	results := make(chan string, len(urls))
	var wg sync.WaitGroup

	// Start the API calls with retries and individual timeouts
	for i, url := range urls {
		ctx, cancel := context.WithTimeout(context.Background(), delays[i])
		defer cancel()
		wg.Add(1)
		go func(url string, delay time.Duration, maxRetries int, ctx context.Context) {
			defer wg.Done()
			retryApiCall(ctx, url, delay, maxRetries, results)
		}(url, delays[i], maxRetries, ctx)