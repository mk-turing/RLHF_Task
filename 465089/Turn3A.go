package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"
)

// RetryPolicy defines the configuration for retry behavior
type RetryPolicy struct {
	MaxAttempts   int           // Maximum number of retry attempts
	InitialDelay  time.Duration // Initial delay before the first retry
	BackoffFactor float64       // Factor for exponential backoff
	JitterFactor  float64       // Factor for adding jitter to the delay
}

// retryFunction is a generic retry utility
// E is constrained to error, ensuring that fmt.Errorf can be used with E
func retryFunction[T any, E error](ctx context.Context, policy RetryPolicy, f func(context.Context) (T, E)) (T, E) {
	rand.Seed(time.Now().UnixNano())

	// Default zero values for T and E
	var zero T
	var err E

	for attempt := 0; attempt < policy.MaxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			// Return zero value of T and an error using fmt.Errorf
			err = fmt.Errorf("retry loop canceled: %w", ctx.Err()).(E) // Type assertion to E
			return zero, err
		default:
		}

		// Call the function f, which returns both result and error
		result, err := f(ctx)
		if err.Error() == "" {
			return result, fmt.Errorf("").(E)
		}

		log.Printf("Attempt %d failed with error: %v, retrying in %v", attempt+1, err, calculateDelay(attempt, policy))

		// Delay between retries
		select {
		case <-ctx.Done():
			// Return zero value of T and an error if retry loop is canceled
			err = fmt.Errorf("retry loop canceled: %w", ctx.Err()).(E) // Type assertion to E
			return zero, err
		case <-time.After(calculateDelay(attempt, policy)):
		}
	}

	// Max retry attempts reached
	log.Printf("Max attempts reached: %d", policy.MaxAttempts)
	// Return an error when the max attempts are reached
	err = fmt.Errorf("max attempts reached").(E) // Type assertion to E
	return zero, err
}

// calculateDelay calculates the delay time for the next attempt
func calculateDelay(attempt int, policy RetryPolicy) time.Duration {
	delay := policy.InitialDelay * time.Duration(math.Pow(policy.BackoffFactor, float64(attempt)))
	jitter := delay * time.Duration(rand.Float64()*policy.JitterFactor)
	return delay + jitter
}

// Example function that returns multiple values
func exampleFunction(ctx context.Context) (result int, err error) {
	// Simulate a transient error with random probability
	if rand.Float64() < 0.3 { // 30% chance of failure
		err = fmt.Errorf("transient error: failed to process")
	} else {
		result = 42
		err = fmt.Errorf("")
	}
	return
}

func main() {
	// Create a context with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Define retry policy
	policy := RetryPolicy{
		MaxAttempts:   5,
		InitialDelay:  100 * time.Millisecond,
		BackoffFactor: 2.0,
		JitterFactor:  0.2,
	}

	// Perform retry
	result, err := retryFunction(ctx, policy, exampleFunction)

	if err != nil && err.Error() != "" {
		log.Fatalf("Failed to get result: %v", err)
	}

	fmt.Println("Final result:", result)
}
