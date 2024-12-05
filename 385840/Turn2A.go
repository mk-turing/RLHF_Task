package main

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"
)

const (
	maxRetries     = 3
	initialBackoff = 100 * time.Millisecond
	maxBackoff     = 10 * time.Second
)

// Simulate a remote operation that might fail transiently
func remoteOperation(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if rand.Intn(2) == 0 {
			return fmt.Errorf("transient failure")
		}
		return nil
	}
}

// Perform retries with exponential backoff
func retryWithExponentialBackoff(ctx context.Context, f func(context.Context) error) error {
	for retry := 0; retry <= maxRetries; retry++ {
		err := f(ctx)
		if err == nil {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		backoff := initialBackoff * time.Duration(math.Pow(2, float64(retry)))
		if backoff > maxBackoff {
			backoff = maxBackoff
		}

		fmt.Printf("Retry %d failed, backing off for %s...\n", retry, backoff)
		time.Sleep(backoff)
	}

	return fmt.Errorf("max retries exceeded")
}

func main() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	defer cancel()

	const numGoroutines = 5

	for i := 0; i < numGoroutines; i++ {
		name := fmt.Sprintf("worker %d", i)
		go func(name string) {
			err := retryWithExponentialBackoff(ctx, func(ctx context.Context) error {
				fmt.Printf("%s: Starting remote operation...\n", name)
				return remoteOperation(ctx)
			})
			if err != nil {
				fmt.Printf("%s: Remote operation failed: %v\n", name, err)
			} else {
				fmt.Printf("%s: Remote operation succeeded\n", name)
			}
		}(name)
	}

	// Simulate some work
	time.Sleep(3 * time.Second)
	rand.Seed(time.Now().UnixNano())
	if rand.Intn(2) == 0 {
		fmt.Println("Cancelling operations")
		cancel()
	}

	select {
	case <-ctx.Done():
		fmt.Println("Main finished")
	}
}
