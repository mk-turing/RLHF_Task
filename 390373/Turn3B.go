package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

// Custom error structure for retryable errors
type RetryableError struct {
	Err        error
	RetryCount int
}

func (re *RetryableError) Error() string {
	return fmt.Sprintf("error: %v, retry count: %d", re.Err, re.RetryCount)
}

// Example error categorization structure
type NonRetryableError struct {
	Err error
}

func (nre *NonRetryableError) Error() string {
	return fmt.Sprintf("non-retryable error: %v", nre.Err)
}

// CircuitBreaker struct to manage the circuit breaker state
type CircuitBreaker struct {
	state             string
	successThreshold int
	totalRequests     int
	successRequests   int
	failureThreshold int
	lastFailureTime   time.Time
	mu                 sync.Mutex
	delay               time.Duration
	closeAfterSuccess  int
}

func NewCircuitBreaker(successThreshold, failureThreshold int, delay time.Duration, closeAfterSuccess int) *CircuitBreaker {
	return &CircuitBreaker{
		state:              "closed",
		successThreshold:   successThreshold,
		failureThreshold:   failureThreshold,
		delay:               delay,
		closeAfterSuccess: closeAfterSuccess,
	}
}

func (cb *CircuitBreaker) isClosed() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state == "closed"
}

func (cb *CircuitBreaker) markSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.successRequests++
	cb.totalRequests++
}

func (cb *CircuitBreaker) markFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.totalRequests++
	cb.failureThreshold++
	cb.lastFailureTime = time.Now()
}

func (cb *CircuitBreaker) tryHalfOpen() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	if time.Since(cb.lastFailureTime) > cb.delay {
		cb.state = "half_open"
		cb.successRequests = 0
		cb.failureThreshold = 0
		return true
	}
	return false
}

func (cb *CircuitBreaker) close() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.state = "closed"
	cb.failureThreshold = 0
}

func (cb *CircuitBreaker) open() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.state = "open"
	log.Println("Circuit breaker is OPEN")
}

// Simulates an operation that may fail
func operation(ctx context.Context) error {
	if rand.Float32() < 0.7 { // ~70% chance to fail
		// Randomly return a type of error
		if rand.Float32() < 0.5 {
			return &RetryableError{
				Err: fmt.Errorf("temporary error occurred"),
			}
		}
		return &NonRetryableError{
			Err: fmt.Errorf("permanent error occurred"),
		}
	}
	return nil // operation succeeded
}

// Circuit-breaked retry logic
func retryOperationWithBreaker(ctx context.Context, maxRetries int, breaker *CircuitBreaker) error {
	retries := 0
	backoff := time.Second // starting backoff duration

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if !breaker.isClosed() {
				log.Println("Circuit breaker is OPEN or HALF_OPEN, skipping operation")
				return fmt.Errorf("service unavailable due to circuit breaker")
			}

			err := operation(ctx)
			if err == nil {
				breaker.markSuccess()
				return nil // success
			}

			// Handle different error types
			if retryableErr, ok := err.(*RetryableError); ok {
				retryableErr.RetryCount = retries // update count
				log.Printf("Retryable error occurred: %s", retryableErr)
				log.Printf("Attempting retry %d after %v", retries+1, backoff)
				retries++
				if retries > maxRetries {
					breaker.markFailure()
					breaker.open()
					return fmt.Errorf("max retries exceeded: %w", retryableErr)
				}
				time.Sleep(backoff) // wait before retrying

				// Exponential backoff
				backoff *= 2 // double the backoff time
				continue // retry the operation
			}

			// If it's a non-retryable error
			if _, ok := err.(*NonRetryableError); ok {
				log.Printf("Non-retryable error encountered: %s", err)
				breaker.markFailure()
				breaker.open()
				return fmt.Errorf("operation failed with non-retryable error: %w", err)
			}

			// Log unexpected error (should not happen)
			log.Printf("Unexpected error: %s", err)
			breaker.markFailure()
			breaker.open()
			return fmt.Errorf("unexpected error occurred: %w", err)
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano()) // Seed random number generator

	cb := NewCircuitBreaker(successThreshold:=3, failureThreshold:=2, delay:=5*time.Second, closeAfterSuccess:=3)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second) // Context with timeout
	defer cancel()

	for {
		if cb.tryHalfOpen() {
			log.Println("Circuit breaker is HALF_OPEN, attempting one request...")
			if err := retryOperationWithBreaker(ctx, 5, cb); err != nil {
				log.Printf("Operation failed: %v, closing circuit breaker", err)
				cb.close()
			} else {
				cb.close()
				log.Println("Operation succeeded, circuit breaker closed")
			}
		}

		if cb.isClosed() {
			if err := retryOperationWithBreaker(ctx, 5, cb); err != nil {
				log.Printf("Operation failed: %v", err)
			} else {
				log.Println("Operation succeeded")
			}
		}

		time.Sleep(2 * time.Second)
	}
}