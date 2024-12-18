package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// simulateApiCall makes an HTTP request and simulates a delay.
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

// handleApiRequest demonstrates how to handle API requests with context.
func handleApiRequest(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	delayStr := r.URL.Query().Get("delay")
	delay, err := time.ParseDuration(delayStr)
	if err != nil {
		http.Error(w, "Invalid delay parameter", http.StatusBadRequest)
		return
	}

	// Create a new context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Perform the API call
	response, err := simulateApiCall(ctx, url, delay)
	if err != nil {
		if ctx.Err() != nil {
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Write the response
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func main() {
	http.HandleFunc("/api", handleApiRequest)
	fmt.Println("Server listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
