package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

// RetryPolicy defines the configuration for the retry mechanism.
type RetryPolicy struct {
	MaxRetries     int         // Maximum number of retry attempts.
	BackoffFactor  float64     // Factor for exponential backoff.
	JitterRatio    float64     // Ratio of random jitter to add to backoff.
	InitialBackoff time.Duration // Initial backoff duration.
}

// RetryFunc is the function type that will be retried.
type RetryFunc func(context.Context) (interface{}, error)

// Retry returns the result of the provided RetryFunc, after retrying in case of errors.
// It adheres to the configured RetryPolicy.
func Retry(ctx context.Context, policy RetryPolicy, retryFunc RetryFunc) (interface{}, error) {
	var result interface{}
	var err error

	for attempt := 0; attempt <= policy.MaxRetries; attempt++ {
		// Check if the context is canceled
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("retry loop canceled: %w", ctx.Err())
		default:
		}

		result, err = retryFunc(ctx)
		if err == nil {
			// Success, break out of the retry loop
			break
		}

		// Handle transient errors only
		if err != nil {
			backoffDuration := getBackoffDuration(attempt, policy)
			log.Printf("Attempt %d failed with error: %v, retrying in %v", attempt+1, err, backoffDuration)
			time.Sleep(backoffDuration)
		}
	}

	if err != nil {
		log.Printf("Max attempts reached: %d", policy.MaxRetries)
		return nil, fmt.Errorf("max attempts reached: %w", err)
	}

	return result, nil
}

func getBackoffDuration(attempt int, policy RetryPolicy) time.Duration {
	// Calculate the backoff duration using exponential backoff with jitter
	backoff := float64(policy.InitialBackoff) * policy.BackoffFactor**float64(attempt)
	jitter := time.Duration(rand.Float64() * policy.JitterRatio * float64(backoff))
	return time.Duration(backoff) + jitter
}

func exampleFunction1(ctx context.Context) (string, error) {
	// Simulate a transient error
	if rand.Float64() < 0.5 {
		return "", fmt.Errorf("transient error in exampleFunction1")
	}
	return "success", nil
}

func exampleFunction2(ctx context.Context, input int) ([]int, error) {
	// Simulate a transient error for input < 3
	if input < 3 {
		return nil, fmt.Errorf("transient error in exampleFunction2: input is too low")
	}
	return []int{input * 2, input * 3}, nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Define retry policies for different functions
	policy1 := RetryPolicy{
		MaxRetries:     3,
		BackoffFactor:  2,
		JitterRatio:    0.2,
		InitialBackoff: 200 * time.Millisecond,
	}

	policy2 := RetryPolicy{
		MaxRetries:     5,
		BackoffFactor:  1.5,
		JitterRatio:    0.3,
		InitialBackoff: 100 * time.Millisecond,