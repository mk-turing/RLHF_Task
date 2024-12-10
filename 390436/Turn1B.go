package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"
)

type RateLimitError struct {
	Code      int       `json:"code" xml:"code"`
	Message   string    `json:"message" xml:"message"`
	Limit     int       `json:"limit" xml:"limit"`
	Remaining int       `json:"remaining" xml:"remaining"`
	Reset     time.Time `json:"reset" xml:"reset"`
}

// Implement the Error() method for RateLimitError
func (e *RateLimitError) Error() string {
	return fmt.Sprintf("RateLimitError: %s (Code: %d)", e.Message, e.Code)
}

func rateLimitHandler(w http.ResponseWriter, r *http.Request) {
	// Simulate checking rate limit
	if checkRateLimit(r.RemoteAddr) {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	// Process request as usual
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello, %s!", r.RemoteAddr)
}

func checkRateLimit(remoteAddr string) bool {
	// Simple in-memory rate limiting simulation
	limit := 10
	interval := time.Minute
	now := time.Now()

	// Replace this with a real rate limiting mechanism
	if now.Sub(lastRequest[remoteAddr]) < interval {
		return requests[remoteAddr] >= limit
	}

	requests[remoteAddr] = 1
	lastRequest[remoteAddr] = now
	return false
}

func handleRateLimitError(w http.ResponseWriter, r *http.Request, err error) {
	rateLimitErr, ok := err.(*RateLimitError)
	if !ok {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rateLimitErr.Code)

	if err := json.NewEncoder(w).Encode(rateLimitErr); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/rate", rateLimitHandler)

	// Add error handling middleware
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				handleRateLimitError(w, r, err.(error))
			}
		}()
		rateLimitHandler(w, r)
	})

	// Example XML handling (add as needed)
	http.HandleFunc("/xml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusTooManyRequests)

		rateLimitErr := &RateLimitError{
			Code:      http.StatusTooManyRequests,
			Message:   "You have exceeded the rate limit.",
			Limit:     10,
			Remaining: 0,
			Reset:     time.Now().Add(time.Minute),
		}

		if err := xml.NewEncoder(w).Encode(rateLimitErr); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	fmt.Println("Server listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

// Simulated rate limiting state
var requests = make(map[string]int)
var lastRequest = make(map[string]time.Time)
