package main

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

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

// Centralized Logger using logrus
var logger = logrus.New()

func init() {
	// Set log format and output
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logger.SetOutput(os.Stderr)
}

// Service Layer
type UserService struct{}

func (s *UserService) GetUser(ctx context.Context, id string) (string, error) {
	if id == "" {
		return "", &ApplicationError{Code: 400, Message: "user ID cannot be empty", Err: nil}
	}

	return "UserData", nil
}

func (s *UserService) ProcessPayment(ctx context.Context) error {
	return &NetworkError{Message: "Payment gateway timeout", Err: fmt.Errorf("timeout")}
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

// Centralized Error Logging
func LogError(ctx context.Context, err error) {
	logger.Error(err)
}

// Main Application Entry Point
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
