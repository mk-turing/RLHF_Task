package main

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// --- Tracing Setup ---
func InitTracer() {
	// Create a simple stdout exporter
	exp, err := stdout.NewExporter(stdout.WithPrettyPrint())
	if err != nil {
		log.Fatalf("failed to create trace exporter: %v", err)
	}

	// Create a new TracerProvider with the exporter
	tp := trace.NewTracerProvider(trace.WithBatcher(exp))
	otel.SetTracerProvider(tp)
}

func StartSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	tracer := otel.Tracer("service-name")
	return tracer.Start(ctx, name)
}

// --- Retry Mechanism ---
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

		// Simulate checking if the error is transient
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
	// Assume network-related errors are transient
	switch err.(type) {
	case *NetworkError:
		return true
	default:
		return false
	}
}

// --- Error Definitions ---
type NetworkError struct {
	Message string
	Err     error
}

func (e *NetworkError) Error() string {
	return e.Message
}

// --- Service Layer ---
type UserService struct{}

func (s *UserService) GetUser(ctx context.Context, id string) (string, error) {
	if id == "" {
		return "", errors.New("user ID cannot be empty")
	}

	return "UserData", nil
}

func (s *UserService) ProcessPayment(ctx context.Context) error {
	ctx, span := StartSpan(ctx, "ProcessPayment")
	defer span.End()

	policy := RetryPolicy{
		MaxAttempts: 3,
		Backoff:     ExponentialBackoff,
	}

	return Retry(ctx, policy, func(ctx context.Context) error {
		if rand.Intn(2) == 0 {
			return &NetworkError{Message: "Payment gateway timeout", Err: errors.New("timeout")}
		}
		return nil
	})
}

func ExponentialBackoff(attempt int) time.Duration {
	base := time.Millisecond * 100 // 100ms
	cap := time.Second * 10        // 10s
	return time.Duration(min(int64(base*time.Pow(2, float64(attempt))), int64(cap)))
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// --- HTTP Handlers ---
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
	traceID := trace.SpanFromContext(ctx).SpanContext().TraceID.String()
	log.Printf("Error with TraceID %s: %v", traceID, err)

	switch err := err.(type) {
	case *NetworkError:
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	default:
		http.Error(w, "Unexpected error", http.StatusInternalServerError)
	}
}

// --- Main Function ---
func main() {
	InitTracer() // Initialize OpenTelemetry tracing

	handler := &Handler{
		userService: UserService{},
	}

	http.HandleFunc("/user", handler.GetUserHandler)
	http.HandleFunc("/payment", handler.PaymentHandler)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
