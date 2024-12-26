package main

import (
	"fmt"
	"net/http"
	"reflect"
	"sync"
	"time"
)

type Logger struct{}

func (l *Logger) Log(msg string) {
	fmt.Println("Log:", msg)
}

type Middleware struct {
	Logger         *Logger
	RateLimiter    RateLimiter
	RequestCounter int64
	MaxRequests    int64
	TimeWindow     time.Duration
}

func (m *Middleware) Handle(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		key := m.getKey(params)

		if !m.RateLimiter.CanRequest(key) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		m.Logger.Log("Query Parameters:")
		for key, values := range params {
			// Extract the first value for logging
			value := values[0]
			m.Logger.Log(fmt.Sprintf("%s: %v (%T)", key, value, reflect.TypeOf(value)))
		}

		m.RateLimiter.Increment(key)
		h.ServeHTTP(w, r)
	})
}

func (m *Middleware) getKey(params http.Values) string {
	var key string
	for _, values := range params {
		key += values[0] + ","
	}
	return key[:len(key)-1]
}

type RateLimiter interface {
	CanRequest(key string) bool
	Increment(key string)
}

type inMemoryRateLimiter struct {
	mu       sync.RWMutex
	cache    map[string]int64
	timeout  time.Duration
	maxCount int64
}

func NewInMemoryRateLimiter(maxCount int64, timeout time.Duration) RateLimiter {
	return &inMemoryRateLimiter{
		cache:    make(map[string]int64),
		timeout:  timeout,
		maxCount: maxCount,
	}
}

func (l *inMemoryRateLimiter) CanRequest(key string) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()

	now := time.Now()
	count, ok := l.cache[key]

	if !ok || now.Sub(time.Unix(count, 0)) > l.timeout {
		return true // Allow request if key is not found or timeout elapsed
	}

	if count >= l.maxCount {
		return false // Exceeded the maximum count for this time window
	}

	return true
}

func (l *inMemoryRateLimiter) Increment(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	count, ok := l.cache[key]

	if !ok {
		l.cache[key] = now.Unix()
		return
	}

	if now.Sub(time.Unix(count, 0)) > l.timeout {
		l.cache[key] = now.Unix()
		return
	}

	l.cache[key] = now.Unix()
}

func main() {
	l := &Logger{}
	limiter := NewInMemoryRateLimiter(10, 1*time.Minute) // Limit 10 requests per minute
	m := &Middleware{Logger: l, RateLimiter: limiter, MaxRequests: 10, TimeWindow: 1 * time.Minute}

	http.Handle("/", m.Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Request processed successfully.")
	})))

	fmt.Println("Server is running on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
