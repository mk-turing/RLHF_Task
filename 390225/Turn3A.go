package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	mu       sync.Mutex
	tokens   int
	lastTime time.Time
	rate     int // tokens per second
}

func NewRateLimiter(rate int) *RateLimiter {
	return &RateLimiter{tokens: rate, lastTime: time.Now(), rate: rate}
}

func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastTime).Seconds()
	rl.tokens += int(elapsed * float64(rl.rate))
	if rl.tokens > rl.rate {
		rl.tokens = rl.rate
	}
	rl.lastTime = now

	if rl.tokens > 0 {
		rl.tokens--
		return true
	}
	return false
}

func (rl *RateLimiter) rateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !rl.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next(w, r)
	}
}

func callServiceWithRetry(url string, retries int) (string, error) {
	client := &http.Client{Timeout: 2 * time.Second}
	var err error
	for i := 0; i < retries; i++ {
		resp, err := client.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			return string(body), nil
		}
		log.Printf("Attempt %d: Error calling service: %v", i+1, err)
		time.Sleep(100 * time.Millisecond)
	}
	return "", fmt.Errorf("failed to call service after %d attempts: %v", retries, err)
}

func errorHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, fmt.Sprintf("Internal Server Error: %v", err), http.StatusInternalServerError)
				log.Printf("Recovered from error: %v", err)
			}
		}()
		next(w, r)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Example services
	urlA := "http://localhost:8081/serviceA"
	urlB := "http://localhost:8082/serviceB"

	var wg sync.WaitGroup
	results := make([]string, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if result, err := callServiceWithRetry(urlA, 3); err == nil {
			results[0] = result
		} else {
			results[0] = fmt.Sprintf("Error: %s", err.Error())
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if result, err := callServiceWithRetry(urlB, 3); err == nil {
			results[1] = result
		} else {
			results[1] = fmt.Sprintf("Error: %s", err.Error())
		}
	}()

	wg.Wait()
	fmt.Fprintln(w, results)
}

func main() {
	rateLimiter := NewRateLimiter(5) // 5 requests per second

	http.HandleFunc("/gateway", rateLimiter.rateLimitMiddleware(errorHandler(handleRequest)))

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
