package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"os"
	"time"
)

// Logger setup using logrus
var logger = logrus.New()

func init() {
	// Set log format and output
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logger.SetOutput(os.Stderr)
}

// Error Types Definitions
type NetworkError struct {
	Message string
	Err     error
}

func (e *NetworkError) Error() string {
	return fmt.Sprintf("NetworkError: %s", e.Message)
}

func (e *NetworkError) IsTransient() bool {
	return true
}

type ApplicationError struct {
	Code    int
	Message string
	Err     error
}

func (e *ApplicationError) Error() string {
	return fmt.Sprintf("ApplicationError: Code: %d, %s", e.Code, e.Message)
}

// Retry Logic (Exponential Backoff)
type RetryPolicy struct {
	MaxAttempts int
	Backoff     func(attempt int) time.Duration
}

func Retry(ctx context.Context, policy RetryPolicy, f func(context.Context) error) error {
	attempt := 0
	for attempt < policy.MaxAttempts {
		attempt++
		err := f(ctx)
		if err == nil {
			return nil
		}

		// Check if the error is transient
		if isTransientError(err) {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(policy.Backoff(attempt)):
			}
		} else {
			return err
		}
	}
	return errors.New("max retries exceeded")
}

func isTransientError(err error) bool {
	// Implement logic to determine if an error is transient
	switch err.(type) {
	case *NetworkError:
		return true
	default:
		return false
	}
}

// Exponential Backoff Strategy
func ExponentialBackoff(attempt int) time.Duration {
	base := time.Millisecond * 100 // 100ms
	//cap := time.Second * 10        // 10s
	return time.Duration(base * time.Second)
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// Service Layer
type UserService struct{}

func (s *UserService) GetUser(ctx context.Context, id string) (string, error) {
	if id == "" {
		return "", errors.New("user ID cannot be empty")
	}

	// Simulating a successful data retrieval
	return "UserData", nil
}

func (s *UserService) ProcessPayment(ctx context.Context) error {
	policy := RetryPolicy{
		MaxAttempts: 3,
		Backoff:     ExponentialBackoff,
	}

	return Retry(ctx, policy, func(ctx context.Context) error {
		// Simulate a transient network error
		if rand.Intn(2) == 0 {
			return &NetworkError{Message: "Payment gateway timeout", Err: errors.New("timeout")}
		}
		return nil
	})
}

// Handler Layer
type Handler struct {
	userService UserService
}

func (h *Handler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), "request", r)
	id := r.URL.Query().Get("id")
	userData, err := h.userService.GetUser(ctx, id)
	if err != nil {
		h.handleError(ctx, w, err)
		return
	}

	w.Write([]byte(userData))
}

func (h *Handler) PaymentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), "request", r)
	err := h.userService.ProcessPayment(ctx)
	if err != nil {
		h.handleError(ctx, w, err)
		return
	}

	w.Write([]byte("Payment successful"))
}

func (h *Handler) handleError(ctx context.Context, w http.ResponseWriter, err error) {
	switch err := err.(type) {
	case *ApplicationError:
		http.Error(w, err.Error(), err.Code)
	case *NetworkError:
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	default:
		http.Error(w, "Unexpected error", http.StatusInternalServerError)
	}
	LogError(ctx, err)
}

// Log Error
func LogError(ctx context.Context, err error) {
	logger.WithContext(ctx).Errorf("Error: %s", err)
}

// Main Function to Start the Server
func main() {
	handler := &Handler{
		userService: UserService{},
	}

	http.HandleFunc("/user", handler.GetUserHandler)
	http.HandleFunc("/payment", handler.PaymentHandler)

	logger.Info("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Fatal(err)
	}
}
