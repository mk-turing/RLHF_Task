package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// RateLimiter struct to manage rate limiting
type RateLimiter struct {
	requests map[string]int
	mu       sync.Mutex
	limit    int
	window   time.Duration
}

// ErrorResponse struct for custom error messages
type ErrorResponse struct {
	Error        string `json:"error"`
	Message      string `json:"message,omitempty"`
	RetryAfter   int    `json:"retry_after,omitempty"`
	CurrentCount int    `json:"current_count,omitempty"`
	Limit        int    `json:"limit,omitempty"`
}

// NewRateLimiter creates a new RateLimiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string]int),
		limit:    limit,
		window:   window,
	}
}

// limitRequest checks if the IP is within the allowed request limit
func (rl *RateLimiter) limitRequest(ip string) (bool, int) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	count := rl.requests[ip]
	if count >= rl.limit {
		return false, count
	}
	rl.requests[ip]++

	// Reset count after the time window
	go func() {
		time.Sleep(rl.window)
		rl.mu.Lock()
		rl.requests[ip]--
		rl.mu.Unlock()
	}()

	return true, count
}

// middleware for rate limiting
func (rl *RateLimiter) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr // retrieve the user's IP (simplified)

		allowed, count := rl.limitRequest(ip)
		if !allowed {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests) // 429 status code

			// Build custom error response
			waitTime := int(rl.window.Seconds())
			response := ErrorResponse{
				Error:        "Rate limit exceeded",
				Message:      "You have exceeded your request limit. Please wait and try again.",
				RetryAfter:   waitTime,
				CurrentCount: count,
				Limit:        rl.limit,
			}

			json.NewEncoder(w).Encode(response) // send JSON response
			return
		}

		next.ServeHTTP(w, r) // pass to next handler if allowed
	})
}

func main() {
	// Initialize the router and rate limiter
	router := mux.NewRouter()
	rateLimiter := NewRateLimiter(5, 10*time.Second) // Persist 5 requests every 10 seconds

	router.Use(rateLimiter.middleware) // Apply the rate limiting middleware

	// Define a simple API resource endpoint
	router.HandleFunc("/api/resource", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Success!"})
	}).Methods("GET")

	// Start the HTTP server
	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", router)
}
