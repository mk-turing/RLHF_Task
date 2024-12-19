package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

// Example function that returns multiple values
func exampleFunction(ctx context.Context, input int) (result int, err error) {
	// Simulate a transient error
	if input < 10 {
		err = fmt.Errorf("transient error: input %d is too low", input)
	} else {
		result = input * input
	}
	return
}

// Retry strategy with exponential backoff
func retryWithExponentialBackoff(ctx context.Context, attempts int, delay time.Duration, f func(context.Context) (int, error)) (int, error) {
	for attempt := 0; attempt < attempts; attempt++ {
		select {
		case <-ctx.Done():
			return 0, fmt.Errorf("retry loop canceled: %w", ctx.Err())
		default:
		}

		result, err := f(ctx)
		if err == nil {
			log.Printf("Attempt %d succeeded: result = %d", attempt+1, result)
			return result, nil
		}

		log.Printf("Attempt %d failed with error: %v, retrying in %v", attempt+1, err, delay)

		// Exponential backoff
		time.Sleep(delay)
		delay *= 2
	}

	log.Printf("Max attempts reached: %d", attempts)
	return 0, fmt.Errorf("max attempts reached")
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	input := 8 // This will trigger a transient error
	result, err := retryWithExponentialBackoff(ctx, 5, 100*time.Millisecond, func(ctx context.Context) (int, error) {
		return exampleFunction(ctx, input)
	})

	if err != nil {
		log.Fatalf("Failed to get result: %v", err)
	}

	fmt.Println("Final result:", result)
}
