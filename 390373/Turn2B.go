package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/rs/zerolog/zerolog"
)

// Custom error structures to categorize errors
type NetworkError struct {
	Err       error
	Retryable bool
}

type TemporaryError struct {
	Err       error
	Retryable bool
}

type PermanentError struct {
	Err       error
	Retryable bool
}

// Implement the Error method for custom error types
func (ne *NetworkError) Error() string {
	return fmt.Sprintf("network error: %v", ne.Err)
}

func (te *TemporaryError) Error() string {
	return fmt.Sprintf("temporary error: %v", te.Err)
}

func (pe *PermanentError) Error() string {
	return fmt.Sprintf("permanent error: %v", pe.Err)
}

// Example function simulating an operation that may fail
func operation(ctx context.Context) error {
	if rand.Float32() < 0.7 { // ~70% chance to fail
		switch rand.Intn(3) {
		case 0:
			return &NetworkError{
				Err:       fmt.Errorf("network timeout"),
				Retryable: true,
			}
		case 1:
			return &TemporaryError{
				Err:       fmt.Errorf("temporary server error"),
				Retryable: true,
			}
		default:
			return &PermanentError{
				Err:       fmt.Errorf("permanent data corruption"),
				Retryable: false,
			}
		}
	}
	return nil
}

// Retry logic with exponential backoff and observability
func retryOperation(ctx context.Context, maxRetries int, initialDelay time.Duration) error {
	logger := zerolog.Ctx(ctx).WithTime().Str("method", "retryOperation")

	retries := 0
	delay := initialDelay
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := operation(ctx)
			if err == nil {
				logger.Info("Operation completed successfully")
				return nil // success
			}

			logger.Warn().Err(err).Msg("Operation failed")

			// Determine retry behavior based on error type
			switch err.(type) {
			case *NetworkError, *TemporaryError:
				if retries < maxRetries {
					logger.Debug().Msgf("Attempting retry (%d) in %v", retries, delay)
					retries++
					time.Sleep(delay) // wait before retrying
					delay *= 2        // exponential backoff
					continue
				}
				logger.Error().Msg("Max retries exceeded")
			default:
				logger.Error().Msg("Permanent error, not retrying")
			}
			return err // exit with error if retries are exhausted or error is non-retryable
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano()) // Seed random number for simulating operation's success/failure

	zerolog.New(os.Stderr).WithLevel(zerolog.DebugLevel).EnrichWith(zerolog.CommonFields()).Install()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Context with timeout
	defer cancel()

	err := retryOperation(ctx, 5, 200*time.Millisecond) // Retry up to 5 times with a 200ms initial delay
	if err != nil {
		log.Fatalf("Operation failed: %v", err)
	}

	log.Println("Operation completed successfully")
}
