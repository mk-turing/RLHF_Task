package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

// simulateApiCall makes an HTTP request with a delay.
func simulateApiCall(ctx context.Context, url string, delay time.Duration) ([]byte, error) {
	// Check if context has been cancelled before starting the request
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Simulate the delay before making the request
	time.Sleep(delay)

	// Create an HTTP client with the given context
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
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

// handleMultipleApiCalls demonstrates how to handle multiple concurrent API calls.
func handleMultipleApiCalls(urls []string, delays []time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()

	for i, url := range urls {
		if i >= len(delays) {
			continue // Handle case where delay list is shorter than url list
		}

		// Create a new context with a timeout
		ctx, cancel := context.WithTimeout(context.Background(), delays[i])
		defer cancel()

		// Perform the API call
		response, err := simulateApiCall(ctx, url, delays[i])
		if err != nil {
			if ctx.Err() != nil {
				fmt.Printf("Request to %s timed out\n", url)
			} else {
				fmt.Printf("Error requesting %s: %v\n", url, err)
			}
			continue
		}

		fmt.Printf("Successfully fetched %s: %s\n", url, string(response))
	}
}

func main() {
	urls := []string{
		"https://example.com/api/data1",
		"https://example.com/api/data2",
		"https://example.com/api/data3",
	}

	// WaitGroup to synchronize the main thread with the goroutines
	var wg sync.WaitGroup

	// Start a goroutine for each API call
	for _, url := range urls {
		wg.Add(1)
		go handleMultipleApiCalls([]string{url}, []time.Duration{time.Second * 1}, &wg)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	fmt.Println("All API calls completed.")
}
