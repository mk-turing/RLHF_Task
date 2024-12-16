package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

// Config structure holding dynamic parameters for the system
type Config struct {
	FailureThreshold int
	SuccessThreshold int
	CooldownPeriod   time.Duration
	MaxRetries       int
	InitialBackoff   time.Duration
}

// ConfigurationService structure maintaining the configuration and notifying clients
type ConfigurationService struct {
	config      Config
	mu          sync.Mutex
	subscribers []chan Config // Simple pub-sub model, using channels for demo purposes
}

// NewConfigurationService initializes the configuration service
func NewConfigurationService() *ConfigurationService {
	return &ConfigurationService{
		config: Config{
			FailureThreshold: 3,
			SuccessThreshold: 2,
			CooldownPeriod:   5 * time.Second,
			MaxRetries:       5,
			InitialBackoff:   1 * time.Second,
		},
	}
}

// UpdateConfig allows for dynamic updating of configuration
func (cs *ConfigurationService) UpdateConfig(newConfig Config) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.config = newConfig
	cs.notifySubscribers(newConfig)
}

// RegisterSubscriber allows a service to subscribe to configuration updates
func (cs *ConfigurationService) RegisterSubscriber() chan Config {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	newChan := make(chan Config)
	cs.subscribers = append(cs.subscribers, newChan)
	newChan <- cs.config // Send current config immediately to the new subscriber
	return newChan
}

// notifySubscribers sends the updated configuration to all subscribers
func (cs *ConfigurationService) notifySubscribers(newConfig Config) {
	for _, subscriber := range cs.subscribers {
		subscriber <- newConfig // Send new config to each subscriber
	}
}

// CircuitBreaker struct to manage its state
type CircuitBreaker struct {
	failureThreshold int
	successThreshold int
	cooldownPeriod   time.Duration

	failureCount    int
	successCount    int
	lastFailureTime time.Time

	state string // "closed", "open", "half-open"
}

// States
const (
	CircuitClosed   = "closed"
	CircuitOpen     = "open"
	CircuitHalfOpen = "half-open"
)

// NewCircuitBreaker initializes a circuit breaker
func NewCircuitBreaker(config Config) *CircuitBreaker {
	return &CircuitBreaker{
		failureThreshold: config.FailureThreshold,
		successThreshold: config.SuccessThreshold,
		cooldownPeriod:   config.CooldownPeriod,
		state:            CircuitClosed,
	}
}

// allow checks if the operation is allowed by the circuit breaker
func (cb *CircuitBreaker) allow() bool {
	if cb.state == CircuitClosed {
		return true
	}
	if cb.state == CircuitOpen && time.Since(cb.lastFailureTime) > cb.cooldownPeriod {
		cb.state = CircuitHalfOpen
		return true
	}
	return cb.state == CircuitHalfOpen
}

// reportSuccess reports a successful attempt
func (cb *CircuitBreaker) reportSuccess() {
	if cb.state == CircuitHalfOpen {
		cb.state = CircuitClosed
		cb.successCount = 0
	}
	cb.successCount++
}

// reportFailure reports a failed attempt
func (cb *CircuitBreaker) reportFailure() {
	cb.failureCount++
	if cb.failureCount >= cb.failureThreshold {
		cb.state = CircuitOpen
		cb.lastFailureTime = time.Now()
		cb.failureCount = 0
	}
}

// Custom error structure for retryable errors
type RetryableError struct {
	Err        error
	RetryCount int
}

func (re *RetryableError) Error() string {
	return fmt.Sprintf("error: %v, retry count: %d", re.Err, re.RetryCount)
}

// Non-retryable error structure
type NonRetryableError struct {
	Err error
}

func (nre *NonRetryableError) Error() string {
	return fmt.Sprintf("non-retryable error: %v", nre.Err)
}

// Simulates an operation that may fail
func operation(ctx context.Context) error {
	if rand.Float32() < 0.7 { // ~70% chance to fail
		if rand.Float32() < 0.5 {
			return &RetryableError{Err: fmt.Errorf("temporary error occurred")}
		}
		return &NonRetryableError{Err: fmt.Errorf("permanent error occurred")}
	}
	return nil // operation succeeded
}

// Exponential Backoff Retry Logic with Circuit Breaker
func retryOperation(ctx context.Context, cb *CircuitBreaker, maxRetries int, backoff time.Duration) error {
	retries := 0
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
				cb.reportSuccess()
				return nil // success
			}

			if retryableErr, ok := err.(*RetryableError); ok {
				retryableErr.RetryCount = retries
				log.Printf("Retryable error occurred: %s", retryableErr)
				log.Printf("Attempting retry %d after %v", retries+1, backoff)
				retries++
				cb.reportFailure()
				if retries > maxRetries {
					return fmt.Errorf("max retries exceeded: %w", retryableErr)
				}
				time.Sleep(backoff)
				backoff *= 2
				continue
			}

			if _, ok := err.(*NonRetryableError); ok {
				cb.reportFailure()
				log.Printf("Non-retryable error encountered: %s", err)
				return fmt.Errorf("operation failed with non-retryable error: %w", err)
			}

			log.Printf("Unexpected error: %s", err)
			return fmt.Errorf("unexpected error occurred: %w", err)
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Initialize Configuration Service
	configService := NewConfigurationService()

	// Register a subscriber for the retry parameters
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	configChan := configService.RegisterSubscriber()

	// Monitoring configuration updates in a separate goroutine
	go func() {
		for config := range configChan {
			log.Printf("Configuration updated: %+v", config)
		}
	}()

	// Start with current configuration
	currentConfig := configService.config
	cb := NewCircuitBreaker(currentConfig)

	// Simulate handling operations and allow for updates
	for {
		if err := retryOperation(ctx, cb, currentConfig.MaxRetries, currentConfig.InitialBackoff); err != nil {
			log.Printf("Operation failed: %v", err)
		} else {
			log.Println("Operation completed successfully")
		}

		// For demo: Dynamically update configuration every 15 seconds
		if time.Since(time.Now().Add(-20*time.Second)) > 15*time.Second {
			newConfig := Config{
				FailureThreshold: 4,
				SuccessThreshold: 3,
				CooldownPeriod:   10 * time.Second,
				MaxRetries:       3,
				InitialBackoff:   2 * time.Second,
			}
			configService.UpdateConfig(newConfig)
			currentConfig = newConfig
			cb = NewCircuitBreaker(currentConfig) // Create new circuit breaker with updated config

			log.Println("Configuration updated for the system.")
		}
		time.Sleep(5 * time.Second) // Avoid busy-waiting
	}
}
