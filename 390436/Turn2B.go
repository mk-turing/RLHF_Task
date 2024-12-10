package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// Error interface for consistent error handling
type Error interface {
	error
	StatusCode() int
}

// RateLimitError struct for rate limit violations
type RateLimitError struct {
	Message      string `json:"error"`
	RetryAfter   int    `json:"retry_after,omitempty"`
	CurrentCount int    `json:"current_count,omitempty"`
	Limit        int    `json:"limit,omitempty"`
	httpStatus   int
}

func (e *RateLimitError) Error() string {
	return e.Message
}

func (e *RateLimitError) StatusCode() int {
	return e.httpStatus
}

func NewRateLimitError(limit int, currentCount int, window time.Duration) *RateLimitError {
	return &RateLimitError{
		Message:      "Rate limit exceeded",
		RetryAfter:   int(window.Seconds()),
		CurrentCount: currentCount,
		Limit:        limit,
		httpStatus:   http.StatusTooManyRequests,
	}
}

type RateLimiter struct {
	requests map[string]int
	mu       sync.Mutex
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string]int),
		limit:    limit,
		window:   window,
	}
}

func (rl *RateLimiter) limitRequest(ip string) Error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	count := rl.requests[ip]
	if count >= rl.limit {
		rl.requests[ip] = count
		return NewRateLimitError(rl.limit, count, rl.window)
	}
	rl.requests[ip]++

	// Reset the count after the time window
	go func() {
		time.Sleep(rl.window)
		rl.mu.Lock()
		rl.requests[ip]--
		rl.mu.Unlock()
	}()

	return nil
}

func handleError(w http.ResponseWriter, r *http.Request, err Error) {
	w.WriteHeader(err.StatusCode())
	response := ErrorResponse{
		Error: err.Error(),
	}

	// Embed specific fields for rate limit errors
	if rateErr, ok := err.(*RateLimitError); ok {
		response.RetryAfter = rateErr.RetryAfter
		response.CurrentCount = rateErr.CurrentCount
		response.Limit = rateErr.Limit
	}

	json.NewEncoder(w).Encode(response)
}

func main() {
	router := mux.NewRouter()
	rateLimiter := NewRateLimiter(5, 10*time.Second) // 5 requests every 10 seconds

	router.HandleFunc("/api/resource", func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		err := rateLimiter.limitRequest(ip)
		if err != nil {
			handleError(w, r, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Success!"})
	}).Methods("GET")

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", router)
}

type ErrorResponse struct {
	Error        string `json:"error"`
	RetryAfter   int    `json:"retry_after,omitempty"`
	CurrentCount int    `json:"current_count,omitempty"`
	Limit        int    `json:"limit,omitempty"`
}
