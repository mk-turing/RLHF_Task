package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// ... (existing code)

func TestRateLimitMiddleware(t *testing.T) {
	// Create a new token bucket with a low capacity for testing
	bucket := NewTokenBucket(1, 1) // 1 token capacity, refill 1 token per second

	// Set up the handler with the rate limiter
	handler := rateLimitMiddleware(bucket)

	// Test case 1: Request within rate limit
	req, _ := http.NewRequest("GET", "/api/resource", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("want status %d, got %d", http.StatusOK, rr.Code)
	}

	// Test case 2: Exceed rate limit
	for i := 0; i < 5; i++ { // Simulate multiple requests
		req, _ := http.NewRequest("GET", "/api/resource", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		// Wait for the bucket to refill before the next request
		time.Sleep(1 * time.Second)
	}

	// Last request should be rate limited
	req, _ = http.NewRequest("GET", "/api/resource", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("want status %d, got %d", http.StatusTooManyRequests, rr.Code)
	}

	//var response ErrorResponse
	//if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
	//	t.Errorf("error decoding JSON response: %v", err)
	//}
	//
	//if response.Error != "Rate limit exceeded, please try again later." {
	//	t.Errorf("want error message '%s', got '%s'", "Rate limit exceeded, please try again later.", response.Error)
	//}
}
