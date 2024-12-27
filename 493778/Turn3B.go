package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type queryValidator struct {
	params       map[string][]validator
	rateLimits   map[string]*rateLimiter
	rateLimitMu sync.Mutex
}

type validator interface {
	Validate(string, *http.Request) error
}

// ... (existing validators remain the same)

type rateLimiter struct {
	limit      int
	burst      int
	lastReset  time.Time
	requests   int
	mu         sync.Mutex
}

func newRateLimiter(limit, burst int) *rateLimiter {
	return &rateLimiter{
		limit:      limit,
		burst:      burst,
		lastReset:  time.Now(),
		requests:   0,
		mu:         sync.Mutex{},
	}
}

func (rl *rateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	if now.After(rl.lastReset.Add(time.Second)) {
		rl.requests = 0
		rl.lastReset = now
	}

	if rl.requests < rl.limit {
		rl.requests++
		return true
	}

	return false
}

func (v *queryValidator) AddRateLimit(param, value string, limit, burst int) {
	v.rateLimitMu.Lock()
	defer v.rateLimitMu.Unlock()

	if v.rateLimits == nil {
		v.rateLimits = make(map[string]*rateLimiter)
	}

	key := fmt.Sprintf("%s:%s", param, value)
	v.rateLimits[key] = newRateLimiter(limit, burst)
}

func (v *queryValidator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for param, validators := range v.params {
			values, ok := r.URL.Query()[param]
			if !ok || len(values) == 0 {
				continue
			}

			for _, value := range values {
				for _, validator := range validators {
					if err := validator.Validate(value, r); err != nil {
						http.Error(w, fmt.Sprintf("Invalid query parameter %s: %v", param, err), http.StatusBadRequest)
						return
					}
				}

				v.rateLimitMu.Lock()
				rl, ok := v.rateLimits[fmt.Sprintf("%s:%s", param, value)]
				v.rateLimitMu.Unlock()

				if ok && !rl.Allow() {
					http.Error(w, fmt.Sprintf("Rate limit exceeded for parameter %s=%s", param, value), http.StatusTooManyRequests)
					return
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	v := newQueryValidator()
	// ... (existing validation rules remain the same)

	// Rate limiting rules
	v.AddRateLimit("action", "create", 10, 5)  // Allow 10 requests per second for "action=create"
	v.AddRateLimit("order", "asc", 5, 3)     // Allow 5 requests per second for "order=asc"

	http.Handle("/api/data", v.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Data endpoint reached")