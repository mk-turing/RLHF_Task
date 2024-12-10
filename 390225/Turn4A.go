package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type CircuitBreaker struct {
	failureThreshold int
	successThreshold int
	failureCount     int
	successCount     int
	state            string
	lastStateChange  time.Time
	timeout          time.Duration
}

func NewCircuitBreaker(failureThreshold int, successThreshold int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		failureThreshold: failureThreshold,
		successThreshold: successThreshold,
		state:            "CLOSED",
		timeout:          timeout,
	}
}

func (cb *CircuitBreaker) Call(fn func() (string, error)) (string, error) {
	switch cb.state {
	case "OPEN":
		if time.Since(cb.lastStateChange) >= cb.timeout {
			cb.state = "HALF_OPEN"
		} else {
			return "", fmt.Errorf("circuit is open")
		}
	case "HALF_OPEN":
		// Allow one request to pass through
		result, err := fn()
		if err != nil {
			cb.failureCount++
			if cb.failureCount > cb.failureThreshold {
				cb.state = "OPEN"
				cb.lastStateChange = time.Now()
			}
			return result, err
		}
		cb.successCount++
		if cb.successCount >= cb.successThreshold {
			cb.state = "CLOSED"
			cb.failureCount = 0
			cb.successCount = 0
		}
		return result, nil
	}

	// In the CLOSED state
	result, err := fn()
	if err != nil {
		cb.failureCount++
		if cb.failureCount > cb.failureThreshold {
			cb.state = "OPEN"
			cb.lastStateChange = time.Now()
		}
		return result, err
	}
	cb.successCount++
	return result, nil
}

func callService(url string) (string, error) {
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func retryWithBackoff(fn func() (string, error), retries int, delay time.Duration) (string, error) {
	for i := 0; i < retries; i++ {
		result, err := fn()
		if err == nil {
			return result, nil
		}
		log.Printf("Retry %d: %s\n", i+1, err)
		time.Sleep(delay)
		delay *= 2 // Exponential backoff
	}
	return "", fmt.Errorf("all retries failed")
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Example service URL
	url := "http://localhost:8081/service"

	circuitBreaker := NewCircuitBreaker(3, 2, 5*time.Second)

	result, err := circuitBreaker.Call(func() (string, error) {
		return retryWithBackoff(func() (string, error) {
			return callService(url)
		}, 3, 100*time.Millisecond)
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, result)
}

func main() {
	http.HandleFunc("/gateway", handleRequest)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
