package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

// CircuitBreaker struct to manage its state
type CircuitBreaker struct {
	failureThreshold int
	successThreshold int
	cooldownPeriod   time.Duration

	failureCount    int
	successCount    int
	lastFailureTime time.Time

	state string // "closed", "open", "half-open"
	mu    sync.Mutex
}

// States
const (
	CircuitClosed   = "closed"
	CircuitOpen     = "open"
	CircuitHalfOpen = "half-open"
)

// NewCircuitBreaker initializes a circuit breaker
func NewCircuitBreaker(failureThreshold, successThreshold int, cooldownPeriod time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		failureThreshold: failureThreshold,
		successThreshold: successThreshold,
		cooldownPeriod:   cooldownPeriod,
		state:            CircuitClosed,
	}
}

// allow checks if the operation is allowed by the circuit breaker
func (cb *CircuitBreaker) allow() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.state == CircuitClosed {
		return true
	}

	if cb.state == CircuitOpen {
		// Check if the cooldown period has passed
		if time.Since(cb.lastFailureTime) > cb.cooldownPeriod {
			cb.state = CircuitHalfOpen
			return true
		}
		return false
	}

	// In Half-Open state, we allow exactly one request to test if the service is back to normal.
	if cb.state == CircuitHalfOpen {
		return true
	}

	return false
}

// reportSuccess reports a successful attempt
func (cb *CircuitBreaker) reportSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.state == CircuitHalfOpen {
		cb.state = CircuitClosed
		cb.successCount = 0 // Reset successful count on success
	}
	cb.successCount++
}

// reportFailure reports a failed attempt and checks if the circuit should open
func (cb *CircuitBreaker) reportFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failureCount++
	if cb.failureCount >= cb.failureThreshold {
		cb.state = CircuitOpen
		cb.lastFailureTime = time.Now()
		cb.failureCount = 0 // Reset the failure count
	}
}

// Custom error structure for retryable errors
type RetryableError struct {
	Err        error
	RetryCount int
}

// Implement the Error method for RetryableError
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

// Exponential Backoff Retry Logic with Circuit Breaker
func retryOperation(ctx context.Context, cb *CircuitBreaker, maxRetries int) error {
	retries := 0
	backoff := time.Second // starting backoff duration

	for {
		if !cb.allow() {
			log.Println("Circuit is open, not allowing operation.")
			time.Sleep(cb.cooldownPeriod) // Wait before next attempt
			continue
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := operation(ctx)
			if err == nil {
				cb.reportSuccess() // Report success to circuit breaker
				return nil         // success
			}

			// Handle different error types
			if retryableErr, ok := err.(*RetryableError); ok {
				retryableErr.RetryCount = retries // update count
				log.Printf("Retryable error occurred: %s", retryableErr)
				log.Printf("Attempting retry %d after %v", retries+1, backoff)
				retries++
				cb.reportFailure() // Report failure to circuit breaker
				if retries > maxRetries {
					return fmt.Errorf("max retries exceeded: %w", retryableErr)
				}
				time.Sleep(backoff) // wait before retrying
				backoff *= 2        // Exponential backoff
				continue            // retry the operation
			}

			// If it's a non-retryable error
			if _, ok := err.(*NonRetryableError); ok {
				cb.reportFailure() // Report failure to circuit breaker
				log.Printf("Non-retryable error encountered: %s", err)
				return fmt.Errorf("operation failed with non-retryable error: %w", err)
			}

			// Log unexpected error (should not happen)
			log.Printf("Unexpected error: %s", err)
			return fmt.Errorf("unexpected error occurred: %w", err)
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano()) // Seed random number generator

	// Initialize Circuit Breaker
	cb := NewCircuitBreaker(3, 2, 5*time.Second) // 3 failures to open, 2 successes to close, and 5 seconds cooldown

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // Context with timeout
	defer cancel()

	err := retryOperation(ctx, cb, 5) // Retry up to 5 times
	if err != nil {
		log.Fatalf("Operation failed: %v", err)
	}

	log.Println("Operation completed successfully")
}
