package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// RateLimiter struct to manage rate limiting with concurrency
type RateLimiter struct {
	limits sync.Map
	limit  int
	window time.Duration
}

// Entry struct to hold the rate limit information for an IP
type Entry struct {
	count int
	ts    time.Time
}

// NewRateLimiter creates a new RateLimiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		limit:  limit,
		window: window,
	}
}

// limitRequest checks if the IP is within the allowed request limit
func (rl *RateLimiter) limitRequest(ip string) bool {
	now := time.Now()

	if entry, ok := rl.limits.LoadOrStore(ip, &Entry{}); ok {
		entry := entry.(*Entry)

		if now.Sub(entry.ts) < rl.window {
			if entry.count >= rl.limit {
				return false // Rate limit exceeded
			}
			entry.count++
		} else {
			entry.count = 1 // Reset count for new window
		}

		entry.ts = now
		rl.limits.Store(ip, entry)
	}

	return true // Allowed
}

// middleware for rate limiting
func (rl *RateLimiter) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr // retrieve the user's IP (simplified)

		if !rl.limitRequest(ip) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests) // 429 status code

			response := map[string]string{"error": "Rate limit exceeded"}
			json.NewEncoder(w).Encode(response) // send JSON response
			return
		}

		next.ServeHTTP(w, r) // pass to next handler if allowed
	})
}

func main() {
	// Initialize the router and rate limiter
	router := mux.NewRouter()
	rateLimiter := NewRateLimiter(5, 10*time.Second) // 5 requests every 10 seconds

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
