package main

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/juju/ratelimit"
)

type Logger struct{}

func (l *Logger) Log(msg string) {
	fmt.Println("Log:", msg)
}

type RateLimiter struct {
	limiter *ratelimit.Bucket
}

func NewRateLimiter(rate float64, burst int) *RateLimiter {
	return &RateLimiter{limiter: ratelimit.NewBucketWithRate(rate, int64(burst))}
}

func (rl *RateLimiter) Allow() bool {
	return rl.limiter.TakeAvailable(1) > 0
}

type Middleware struct {
	Logger      *Logger
	RateLimiter *RateLimiter
}

func (m *Middleware) Handle(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Rate limit
		if !m.RateLimiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		params := r.URL.Query()
		m.Logger.Log("Query Parameters:")

		for key, values := range params {
			// Extract the first value for logging
			value := values[0]
			m.Logger.Log(fmt.Sprintf("%s: %v (%T)", key, value, reflect.TypeOf(value)))
		}

		h.ServeHTTP(w, r)
	})
}

func main() {
	l := &Logger{}
	// Set the desired rate limit per second
	rateLimiter := NewRateLimiter(1.0, 5) // Allow 1 request per second, burst 5 requests
	m := &Middleware{Logger: l, RateLimiter: rateLimiter}

	http.Handle("/", m.Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Request processed successfully.")
	})))

	fmt.Println("Server is running on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
