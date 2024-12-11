package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
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

// Error Interfaces
type Error interface {
	fmt.Stringer
	Log(context.Context)
}

// SystemError Definition
type SystemError struct {
	Message string
	Err     error
}

func (e *SystemError) Error() string {
	return fmt.Sprintf("SystemError: %s", e.Message)
}

func (e *SystemError) Log(ctx context.Context) {
	logger.WithContext(ctx).Errorf("SystemError: %s", e)
}

// NetworkError Definition
type NetworkError struct {
	Message string
	Err     error
}

func (e *NetworkError) Error() string {
	return fmt.Sprintf("NetworkError: %s", e.Message)
}

func (e *NetworkError) Log(ctx context.Context) {
	logger.WithContext(ctx).Errorf("NetworkError: %s", e)
}

func (e *NetworkError) IsTransient() bool {
	return true
}

// ApplicationError Definition
type ApplicationError struct {
	Code    int
	Message string
	Err     error
}

func (e *ApplicationError) Error() string {
	return fmt.Sprintf("ApplicationError: Code: %d, %s", e.Code, e.Message)
}

func (e *ApplicationError) Log(ctx context.Context) {
	logger.WithContext(ctx).Errorf("ApplicationError: Code: %d, %s", e.Code, e)
}

// Centralized Error Logging
func LogError(ctx context.Context, err error) {
	logger.Error(err)
}

func Retry(ctx context.Context, attempts int, backoff time.Duration, operation func(ctx context.Context) error) error {
	var err error
	for i := 0; i < attempts; i++ {
		if err = operation(ctx); err == nil {
			return nil // Success
		}

		// Check if the error is transient before retrying
		if transientErr, ok := err.(interface{ IsTransient() bool }); ok && transientErr.IsTransient() {
			// Log the error and wait before retrying
			time.Sleep(backoff)
			backoff *= 2 // Exponential backoff
			continue
		}

		// If it's not a transient error, break and return
		break
	}
	return err
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
	return &NetworkError{Message: "Payment gateway timeout", Err: errors.New("timeout")}
}

func (s *UserService) ProcessPaymentWithRetry(ctx context.Context) error {
	return Retry(ctx, 3, 100*time.Millisecond, func(ctx context.Context) error {
		return s.ProcessPayment(ctx)
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
	err := h.userService.ProcessPaymentWithRetry(ctx)
	if err != nil {
		h.handleError(ctx, w, err)
		return
	}

	w.Write([]byte("Payment successful"))
}

func (h *Handler) handleError(ctx context.Context, w http.ResponseWriter, err error) {
	switch err := err.(type) {
	case *SystemError:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	case *NetworkError:
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	case *ApplicationError:
		http.Error(w, err.Error(), err.Code)
	default:
		http.Error(w, "Unexpected error", http.StatusInternalServerError)
	}
	LogError(ctx, err)
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
