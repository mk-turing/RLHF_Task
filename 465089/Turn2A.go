package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"
)

// Example function that returns multiple values
func exampleFunction(ctx context.Context, input int) (result int, err error) {
	// Simulate a transient error with random probability
	if rand.Float64() < 0.3 { // 30% chance of failure
		err = fmt.Errorf("transient error: input %d caused a failure", input)
	} else {
		result = input * input
	}
	return
}

// Retry strategy with exponential backoff and context-based cancellation
func retryWithExponentialBackoff(ctx context.Context, attempts int, delay time.Duration) (int, error) {
	for attempt := 0; attempt < attempts; attempt++ {
		select {
		case <-ctx.Done():
			return 0, fmt.Errorf("retry loop canceled: %w", ctx.Err())
		default:
		}

		result, err := exampleFunction(ctx, attempt+1)
		if err == nil {
			log.Printf("Attempt %d succeeded: result = %d", attempt+1, result)
			return result, nil
		}

		log.Printf("Attempt %d failed with error: %v, retrying in %v", attempt+1, err, delay)

		// Exponential backoff with a jitter
		sleepTime := delay + time.Duration(rand.Intn(int(delay.Seconds())))
		time.Sleep(sleepTime)
		delay *= 2
	}

	log.Printf("Max attempts reached: %d", attempts)
	return 0, fmt.Errorf("max attempts reached")
}

func main() {
	// Create a context with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Define retry strategy parameters
	attempts := 5
	delay := 100 * time.Millisecond

	// Perform retry
	result, err := retryWithExponentialBackoff(ctx, attempts, delay)

	if err != nil {
		log.Fatalf("Failed to get result: %v", err)
	}

	fmt.Println("Final result:", result)
}
