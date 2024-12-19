package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"
)

// Example function that returns multiple values
func exampleFunction(ctx context.Context, input int) (result int, err error) {
	// Simulate a transient error with a random chance
	if rand.Intn(10) < 3 {
		err = errors.New("transient error")
	} else {
		result = input * input
	}
	return
}

// Retry strategy with exponential backoff and context-based cancellation
func retryWithExponentialBackoff(ctx context.Context, attempts int, initialDelay time.Duration, f func(context.Context) (int, error)) (int, error) {
	var lastErr error
	for attempt := 0; attempt < attempts; attempt++ {
		// Check if the context has been canceled
		if ctx.Err() != nil {
			return 0, fmt.Errorf("retry loop canceled: %w", ctx.Err())
		}

		result, err := f(ctx)
		if err == nil {
			log.Printf("Attempt %d succeeded: result = %d", attempt+1, result)
			return result, nil
		}

		lastErr = err
		log.Printf("Attempt %d failed with error: %v, retrying in %v", attempt+1, err, initialDelay)

		// Exponential backoff
		time.Sleep(initialDelay)
		initialDelay *= 2
	}

	log.Printf("Max attempts reached: %d", attempts)
	return 0, fmt.Errorf("max attempts reached: %w", lastErr)
}

func main() {
	// Create a context with a 5-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	input := 5
	result, err := retryWithExponentialBackoff(ctx, 3, 200*time.Millisecond, func(ctx context.Context) (int, error) {
		return exampleFunction(ctx, input)
	})

	if err != nil {
		if errors.As(err, ctx.Err()) {
			log.Fatalf("Failed to get result within timeout: %v", err)
		} else {
			log.Fatalf("Failed to get result: %v", err)
		}
	}

	fmt.Println("Final result:", result)
}
