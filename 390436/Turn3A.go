package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type TokenBucket struct {
	capacity   int        // maximum number of tokens
	tokens     int        // current number of tokens
	refillRate int        // tokens to refill per second
	lastRefill time.Time  // last refill time
	mu         sync.Mutex // mutex for concurrency
}

// NewTokenBucket creates a new TokenBucket
func NewTokenBucket(capacity int, refillRate int) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     capacity,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// GetToken tries to get a token from the bucket
func (tb *TokenBucket) GetToken() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	// Refill tokens based on elapsed time
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)
	tokensToAdd := int(elapsed.Seconds()) * tb.refillRate
	if tokensToAdd > 0 {
		tb.tokens += tokensToAdd
		if tb.tokens > tb.capacity {
			tb.tokens = tb.capacity
		}
		tb.lastRefill = now
	}

	// If there's at least one token, consume it
	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	return false
}

func rateLimitMiddleware(tb *TokenBucket) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !tb.GetToken() {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"error": "Rate limit exceeded, please try again later."}`))
			return
		}
		// If we have a token, proceed with the request
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Request successful!"}`))
	})
}

func main() {
	bucket := NewTokenBucket(5, 1) // 5 tokens capacity, refill 1 token per second

	http.Handle("/api/resource", rateLimitMiddleware(bucket))

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
