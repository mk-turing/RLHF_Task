package main

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"
)

// mockServe simulates a server that processes requests with a delay.
func mockServe(delay time.Duration, c chan string) {
	defer close(c)
	for req := range c {
		fmt.Printf("Processing request: %v\n", req)
		time.Sleep(delay)
		c <- fmt.Sprintf("Response for: %v\n", req)
	}
}

// serveRequests simulates sending requests asynchronously to a mock server.
func serveRequests(requests []string, responses chan string) {
	var wg sync.WaitGroup

	for _, req := range requests {
		wg.Add(1)
		go func() {
			fmt.Printf("Sending request: %v\n", req)
			responses <- "Response for: " + req
			wg.Done()
		}()
	}
	wg.Wait()
	close(responses)
}

func TestAsyncAPIProcessing(t *testing.T) {
	// Simulate incoming requests.
	requests := []string{"req1", "req2", "req3"}
	var responses []string

	// Simulate a response channel for mock API responses.
	responsesCh := make(chan string)

	// We want to await responses using another channel, kept separate from inputs.
	responsesChan := make(chan string)

	go mockServe(200*time.Millisecond, responsesCh)
	defer close(responsesCh)

	serveRequests(requests, responsesChan)

	// Wait for all responses.
	for resp := range responsesChan {
		responses = append(responses, resp)
	}

	// Check if all responses have been processed in order.
	expectedResponses := []string{"Response for: req1\n", "Response for: req2\n", "Response for: req3\n"}
	if !reflect.DeepEqual(responses, expectedResponses) {
		t.Fatalf("Expected responses %v, but got %v", expectedResponses, responses)
	}

	t.Logf("All requests processed successfully.")
}
