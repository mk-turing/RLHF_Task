package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/juju/ratelimit"
)

// RateLimiter holds rate limit configuration
type RateLimiter struct {
	limit    *ratelimit.Bucket
	burst    int
	interval time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rps int, burst int) *RateLimiter {
	interval := time.Second / time.Duration(rps)
	limit := ratelimit.NewBucketWithRate(float64(rps), int64(burst))
	return &RateLimiter{limit: limit, burst: burst, interval: interval}
}

// Allow checks if the request is allowed within the rate limit
func (rl *RateLimiter) Allow() bool {
	return rl.limit.TakeAvailable(1) > 0
}

// RateLimitMiddleware is a middleware to handle rate limiting
func RateLimitMiddleware(rl *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !rl.Allow() {
				w.WriteHeader(http.StatusTooManyRequests)
				fmt.Fprintf(w, "Rate limit exceeded. Try again later.")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func main() {
	// Default rate limit configuration
	defaultRPS := 5
	defaultBurst := 10

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, world!")
	})

	http.HandleFunc("/rate-limited", func(w http.ResponseWriter, r *http.Request) {
		// Parse rate limit configuration from URL query parameters
		rpsStr := r.URL.Query().Get("rps")
		burstStr := r.URL.Query().Get("burst")

		rps := defaultRPS
		burst := defaultBurst
		if rpsStr != "" {
			var err error
			rps, err = strconv.Atoi(rpsStr)
			if err != nil {
				log.Printf("Invalid rps parameter: %v. Using default %d", err, defaultRPS)
			}
		}
		if burstStr != "" {
			var err error
			burst, err = strconv.Atoi(burstStr)
			if err != nil {
				log.Printf("Invalid burst parameter: %v. Using default %d", err, defaultBurst)
			}
		}

		rl := NewRateLimiter(rps, burst)
		handler := RateLimitMiddleware(rl)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Rate limited endpoint.")
		}))
		handler.ServeHTTP(w, r)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
