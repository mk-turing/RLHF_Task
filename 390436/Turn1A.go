package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type RateLimiter struct {
	requests map[string]int
	mu       sync.Mutex
	limit    int
	window   time.Duration
}

type ErrorResponse struct {
	Error        string `json:"error"`
	RetryAfter   int    `json:"retry_after,omitempty"`
	CurrentCount int    `json:"current_count,omitempty"`
	Limit        int    `json:"limit,omitempty"`
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string]int),
		limit:    limit,
		window:   window,
	}
}

func (rl *RateLimiter) limitRequest(ip string) (bool, int) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	count := rl.requests[ip]
	if count >= rl.limit {
		return false, count
	}
	rl.requests[ip]++

	// Reset the count after the time window
	go func() {
		time.Sleep(rl.window)
		rl.mu.Lock()
		rl.requests[ip]--
		rl.mu.Unlock()
	}()

	return true, count
}

func (rl *RateLimiter) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr // get user's IP (simple version)

		allowed, count := rl.limitRequest(ip)
		if !allowed {
			w.WriteHeader(http.StatusTooManyRequests)
			// Send custom error response
			waitTime := int(rl.window.Seconds())
			response := ErrorResponse{
				Error:        "Rate limit exceeded",
				RetryAfter:   waitTime,
				CurrentCount: count,
				Limit:        rl.limit,
			}

			json.NewEncoder(w).Encode(response)
			return
		}

		next.ServeHTTP(w, r) // forward to the next handler
	})
}

func main() {
	router := mux.NewRouter()
	rateLimiter := NewRateLimiter(5, 10*time.Second) // 5 requests every 10 seconds

	router.Use(rateLimiter.middleware)

	router.HandleFunc("/api/resource", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Success!"})
	}).Methods("GET")

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", router)
}
