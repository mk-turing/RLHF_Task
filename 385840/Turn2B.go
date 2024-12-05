package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// RetryableFunction represents a function that can be retried.
type RetryableFunction func(context.Context) error

// Retry retries the given function with exponential backoff.
func Retry(ctx context.Context, f RetryableFunction, maxRetries int, initialBackoff time.Duration) error {
	for retry := 0; retry <= maxRetries; retry++ {
		err := f(ctx)
		if err == nil {
			return nil // Success, no need to retry
		}

		select {
		case <-ctx.Done():
			return ctx.Err() // Cancelled
		default:
		}

		// Calculate backoff time
		backoff := time.Duration(rand.Intn(int(initialBackoff.Seconds()))) * time.Second
		fmt.Printf("Retry %d failed with %v, backing off for %v\n", retry+1, err, backoff)
		time.Sleep(backoff)
	}

	return fmt.Errorf("max retries exceeded: %v", maxRetries)
}

func someOperation(ctx context.Context) error {
	// Simulate a transient failure with a 50% chance
	if rand.Intn(2) == 0 {
		return fmt.Errorf("transient failure")
	}
	return nil
}

func main() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	defer cancel()

	if err := Retry(ctx, someOperation, 5, 1*time.Second); err != nil {
		fmt.Println("Operation failed:", err)
	} else {
		fmt.Println("Operation succeeded")
	}
}
