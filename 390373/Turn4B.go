package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/opentracing/opentracing-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Configuration service for dynamic settings
type ConfigService struct {
	mutex            sync.Mutex
	maxRetries       int32
	failureThreshold int32
	successThreshold int32
	cooldownPeriod   time.Duration
	initialBackoff   time.Duration
}

func NewConfigService() *ConfigService {
	return &ConfigService{
		maxRetries:       5,
		failureThreshold: 3,
		successThreshold: 2,
		cooldownPeriod:   5 * time.Second,
		initialBackoff:   1 * time.Second,
	}
}

// Get the current configuration values
func (cs *ConfigService) GetConfig() (maxRetries int, failureThreshold int, successThreshold int, cooldownPeriod time.Duration, initialBackoff time.Duration) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	return int(cs.maxRetries), int(cs.failureThreshold), int(cs.successThreshold), cs.cooldownPeriod, cs.initialBackoff
}

// Update configuration dynamically
func (cs *ConfigService) UpdateConfig(maxRetries int, failureThreshold int, successThreshold int, cooldownPeriod time.Duration, initialBackoff time.Duration) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	atomic.StoreInt32(&cs.maxRetries, int32(maxRetries))
	atomic.StoreInt32(&cs.failureThreshold, int32(failureThreshold))
	atomic.StoreInt32(&cs.successThreshold, int32(successThreshold))
	cs.cooldownPeriod = cooldownPeriod
	cs.initialBackoff = initialBackoff
}

// CircuitBreaker struct to manage its state
type CircuitBreaker struct {
	cs              *ConfigService
	failureCount    int
	successCount    int
	lastFailureTime time.Time
	state           string // "closed", "open", "half-open"
	mu              sync.Mutex
}

// States
const (
	CircuitClosed   = "closed"
	CircuitOpen     = "open"
	CircuitHalfOpen = "half-open"
)

// NewCircuitBreaker initializes a circuit breaker
func NewCircuitBreaker(cs *ConfigService) *CircuitBreaker {
	return &CircuitBreaker{
		cs:    cs,
		state: CircuitClosed,
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
		if time.Since(cb.lastFailureTime) > cb.cs.cooldownPeriod {
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
	if cb.failureCount >= int(cb.cs.failureThreshold) {
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
func operation(ctx context.Context, tracer opentracing.Tracer) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "operation")
	defer span.Finish()

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
	span.SetTag("error", false)
	return nil // operation succeeded
}

// Exponential Backoff Retry Logic with Circuit Breaker
func retryOperation(ctx context.Context, cb *CircuitBreaker, cs *ConfigService) error {
	retries := 0
	backoff := cs.initialBackoff // starting backoff duration

	for {
		if !cb.allow() {
			log.Println("Circuit is open, not allowing operation.")
			time.Sleep(cs.cooldownPeriod) // Wait before next attempt
			continue
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := operation(ctx, otel.Tracer("", nil))
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
				if retries >= int(cs.maxRetries) {
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

	// Initialize Configuration Service
	cs := NewConfigService()

	// Initialize Circuit Breaker
	cb := NewCircuitBreaker(cs)

	// Setup OpenTelemetry Tracing
	otel.SetTracerProvider(sdktrace.NewProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(resource.NewWithAttributes(
			attribute.String("service.name", "retry-circuit-breaker"),
		)),
		sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(otlptrace.NewClient("localhost:4317"))),
	))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // Context with timeout
	defer cancel()

	err := retryOperation(ctx, cb, cs)
	if err != nil {
		log.Fatalf("Operation failed: %v", err)
	}

	log.Println("Operation completed successfully")
}
