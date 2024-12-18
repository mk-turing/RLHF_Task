package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

const (
	maxRetries = 3
	retryDelay = 1 * time.Second
)

// simulateApiCall makes an HTTP request and simulates a delay.
func simulateApiCall(ctx context.Context, url string, delay time.Duration) ([]byte, error) {
	// Check if context has been cancelled before starting the request
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Simulate the delay before making the request
	time.Sleep(delay)

	// Randomly fail the call sometimes
	if rand.Intn(10) == 0 {
		return nil, fmt.Errorf("random failure for %s", url)
	}

	// Create an HTTP client with the given context
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code %d: %s", resp.StatusCode, resp.Status)
	}

	return ioutil.ReadAll(resp.Body)
}

// RetryCall attempts an API call a specified number of times with retry delay.
func RetryCall(ctx context.Context, url string, delay time.Duration, maxRetries int) ([]byte, error) {
	for i := 0; i <= maxRetries; i++ {
		result, err := simulateApiCall(ctx, url, delay)
		if err == nil {
			return result, nil
		}

		// If the context has been cancelled, return the error
		if ctx.Err() != nil {
			return nil, err
		}

		fmt.Printf("Call to %s failed (retry %d of %d): %v\n", url, i+1, maxRetries+1, err)
		time.Sleep(retryDelay)
	}

	return nil, fmt.Errorf("maximum retries exceeded for %s", url)
}

func main() {
	urls := []string{
		"https://example.com/api1",
		"https://example.com/api2",
		"https://example.com/api3",
	}

	delays := []time.Duration{
		2 * time.Second,
		3 * time.Second,
		1 * time.Second,
	}

	results := make(chan string, len(urls))
	var wg sync.WaitGroup

	// Start the API calls with individual contexts and timeouts
	for i, url := range urls {
		ctx, cancel := context.WithTimeout(context.Background(), delays[i])
		defer cancel()

		wg.Add(1)
		go func(url string, delay time.Duration, ctx context.Context) {
			defer wg.Done()

			_, err := RetryCall(ctx, url, delay, maxRetries)
			if err != nil {
				results <- fmt.Sprintf("Error: %v", err)
				return
			}

			results <- fmt.Sprintf("Call to %s succeeded", url)
		}(url, delays[i], ctx)
	}

	wg.Wait()
	close(results)

	// Print the results of the API calls
	for result := range results {
		fmt.Println(result)
	}
}
