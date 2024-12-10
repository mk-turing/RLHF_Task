package main

import (
	"encoding/json"
	main2 "main"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestRateLimit tests the rate limit behavior of the API.
func TestRateLimit(t *testing.T) {
	// Initialize a rate limiter with 2 tokens capacity and 1 token refill per second
	bucket := main2.NewTokenBucket(2, 1)

	// Create a test HTTP handler to test rate limiting
	handler := main2.rateLimitMiddleware(bucket)

	// Create a buffered response writer to capture response
	response := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/resource", nil)

	// First request (should succeed)
	handler.ServeHTTP(response, req)
	if response.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", response.Code)
	}

	// Check the response body
	var body map[string]string
	json.NewDecoder(response.Body).Decode(&body)
	if body["message"] != "Request successful!" {
		t.Errorf("expected success message; got %v", body["message"])
	}

	// Second request (should succeed)
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	if response.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", response.Code)
	}

	// Third request (should exceed limit and return error)
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	if response.Code != http.StatusTooManyRequests {
		t.Errorf("expected status Too Many Requests; got %v", response.Code)
	}

	// Check error message
	var errorBody map[string]string
	json.NewDecoder(response.Body).Decode(&errorBody)
	if errorBody["error"] != "Rate limit exceeded, please try again later." {
		t.Errorf("expected error message; got %v", errorBody["error"])
	}

	// Wait for a while to allow token refill
	time.Sleep(1 * time.Second)

	// Fourth request (should succeed again after refill)
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	if response.Code != http.StatusOK {
		t.Errorf("expected status OK after refill; got %v", response.Code)
	}
}
