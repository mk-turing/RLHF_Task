package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func main() {
	// Create a context with a 5-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Make an API call using the context
	response, err := makeAPICall(ctx)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Process the API response
	fmt.Println("API Response:", response)
}

func makeAPICall(ctx context.Context) (string, error) {
	// Simulate an API call that takes time
	time.Sleep(2 * time.Second)

	// Check if the context has expired
	select {
	case <-ctx.Done():
		return "", errors.New("API call timed out")
	default:
		// Context is not expired, proceed with the API call
	}

	// Perform the actual API call
	// For demonstration purposes, we'll use a simple HTTP request
	req, err := http.NewRequest("GET", "https://example.com", nil)
	if err != nil {
		return "", err
	}

	// Add the context to the HTTP request
	req = req.WithContext(ctx)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body and return it
	// (For simplicity, we'll just return a fixed string in this example)
	return "Hello, World!", nil
}
