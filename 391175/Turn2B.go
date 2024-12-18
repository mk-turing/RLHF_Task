package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// simulateApiCall makes an HTTP request and simulates a delay.
func simulateApiCall(ctx context.Context, url string, delay time.Duration, result chan<- string) {
	// Check if the context has been cancelled before starting the request
	if ctx.Err() != nil {
		result <- fmt.Sprintf("Call to %s cancelled", url)
		return
	}
	// Simulate the delay before making the request
	time.Sleep(delay)
	// Create an HTTP client with the given context
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		result <- fmt.Sprintf("Error creating request for %s: %v", url, err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		result <- fmt.Sprintf("Error sending request for %s: %v", url, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		result <- fmt.Sprintf("Status code %d for %s", resp.StatusCode, url)
		return
	}
	result <- fmt.Sprintf("Call to %s succeeded", url)
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
			simulateApiCall(ctx, url, delay, results)
		}(url, delays[i], ctx)
	}
	wg.Wait()
	close(results)
	// Print the results of the API calls
	for result := range results {
		fmt.Println(result)
	}
}
