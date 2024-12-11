package main

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

// --- Tracing Setup ---
var tracer otel.Tracer

func initTracer() {
	exporter, err := otlptrace.New(otlptrace.WithEndpoint("http://localhost:4317"))
	if err != nil {
		log.Fatalf("failed to create exporter: %v", err)
	}
	otel.SetTracerProvider(trace.NewSDKTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithSyncer(exporter),
		trace.WithResource(resource.NewWithAttributes(
			attribute.String("service.name", "user-payment-service"),
		)),
	))
	tracer = otel.Tracer("user-payment-service")
}

func StartSpan(ctx context.Context, name string) (context.Context, otel.Span) {
	ctx, span := tracer.Start(ctx, name)
	return ctx, span
}

// --- Error Handling ---
type NetworkError struct {
	Message string
	Err     error
}

func (e *NetworkError) Error() string {
	return e.Message
}

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
	switch err.(type) {
	case *NetworkError:
		return true
	default:
		return false
	}
}

// --- Retry Strategy ---
func ExponentialBackoff(attempt int) time.Duration {
	base := time.Millisecond * 100
	cap := time.Second * 10
	return time.Duration(min(int64(base*time.Pow(2, float64(attempt))), int64(cap)))
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// --- Services ---
type UserService struct{}

func (s *UserService) GetUser(ctx context.Context, id string) (string, error) {
	ctx, span := StartSpan(ctx, "GetUser")
	defer span.End()

	if id == "" {
		span.SetStatus(codes.Error, "user ID cannot be empty")
		return "", errors.New("user ID cannot be empty")
	}

	span.SetStatus(codes.Ok, "user data retrieved")
	return "UserData", nil
}

type PaymentService struct{}

func (s *PaymentService) ProcessPayment(ctx context.Context) error {
	ctx, span := StartSpan(ctx, "ProcessPayment")
	defer span.End()

	policy := RetryPolicy{
		MaxAttempts: 3,
		Backoff:     ExponentialBackoff,
	}

	return Retry(ctx, policy, func(ctx context.Context) error {
		if rand.Intn(2) == 0 {
			span.SetStatus(codes.Error, "payment gateway timeout")
			return &NetworkError{Message: "Payment gateway timeout", Err: errors.New("timeout")}
		}
		span.SetStatus(codes.Ok, "payment processed")
		return nil
	})
}

// --- HTTP Handlers ---
type Handler struct {
	userService    UserService
	paymentService PaymentService
}

func (h *Handler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.URL.Query().Get("id")
	userData, err := h.userService.GetUser(ctx, id)
	if err != nil {
		h.handleError(ctx, w, err)
		return
	}

	w.Write([]byte(userData))
}

func (h *Handler) PaymentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := h.paymentService.ProcessPayment(ctx)
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
	initTracer() // Initialize OpenTelemetry tracing

	handler := &Handler{
		userService:    UserService{},
		paymentService: PaymentService{},
	}

	http.HandleFunc("/user", handler.GetUserHandler)
	http.HandleFunc("/payment", handler.PaymentHandler)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
