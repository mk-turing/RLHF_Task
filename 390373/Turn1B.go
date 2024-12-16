package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttracer"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace/status"
	"go.uber.org/zerolog"
	"go.uber.org/zerolog/zap"
)

// Custom error type to encapsulate contextual details
type RetryableError struct {
	err      error
	retriable bool
}

func (e *RetryableError) Error() string {
	return e.err.Error()
}

func (e *RetryableError) IsRetryable() bool {
	return e.retriable
}

// Setup logging with Zerolog
func setupLogger() *zerolog.Logger {
	logger := zerolog.New(os.Stdout).
		Output(zap.ConsoleWriter(zap.UseCaller(), zap.UseLevelEncoder(), zap.UseTimeEncoder()))
	return logger.With().Timestamp().Str("service", "my-service").Logger()
}

// Setup tracing with OpenTelemetry
func setupTracer() (trace.TracerProvider, func(), error) {
	exporter, err := stdouttracer.NewExporter(stdouttracer.Options{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create stdout exporter: %w", err)
	}

	tp := trace.NewSDKTracerProvider(
		trace.WithSynchronousExport(),
		trace.WithResource(resource.NewWithAttributes(
			attribute.String("service.name", "my-service"),
		)),
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithSpanProcessor(trace.NewExportSpanProcessor(exporter)),
	)

	if err := otel.SetTracerProvider(tp); err != nil {
		return nil, nil, fmt.Errorf("failed to set tracer provider: %w", err)
	}

	return tp, func() {
		if err := exporter.Shutdown(context.Background()); err != nil {
			log.Fatalf("failed to shut down exporter: %v", err)
		}
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatalf("failed to shut down tracer provider: %v", err)
		}
	}, nil
}

// Simulate a service call that may fail
func simulateServiceCall(ctx context.Context) error {
	if rand.Intn(2) == 0 {
		return &RetryableError{err: fmt.Errorf("service call failed: transient error"), retriable: true}
	}
	return nil
}

// Retry logic with context-aware retries
func retryWithLoggingAndTracing(ctx context.Context, f func(context.Context) error, maxRetries int, backoff time.Duration) error {
	tracer := otel.Tracer("my-service")
	span, err := tracer.Start(ctx, "retryWithLoggingAndTracing")
	if err != nil {
		log.Err(err).Msg("failed to start retry span")
		return err
	}
	defer span.End()

	for retryCount := 0; retryCount <= maxRetries; retryCount++ {
		log.Ctx(ctx).Info().Str("retry_count", fmt.Sprint(retryCount)).Msg("attempting service call")

		err := f(ctx)
		if err == nil {
			log.Ctx(ctx).Info().Msg("service call successful")
			span.SetStatus(status.New(codes.OK, "service call successful"))
			return nil
		}

		retryableErr, ok := err.(*RetryableError)
		if !ok || !retryableErr.IsRetryable() {
			log.Ctx(ctx).Err(err).Msg("service call failed: non-retryable error")
			span.SetStatus(status.New(codes.Unavailable, err.Error()))
			return err
		}

		log.Ctx(ctx).Warn().Err(err).Msg("service call failed: retrying")
		span.SetStatus(status.New(codes.Unavailable, err.Error()))

		time.Sleep(backoff)
		backoff *= 2 // Example: exponential backoff
	}

	log.Ctx(ctx).Err(fmt.Errorf("exceeded max retries: %d", maxRetries)).Msg("service call failed")
	span.SetStatus(status.New(codes.Unavailable, fmt.Sprintf("exceeded max retries: %d", maxRetries)))
	return fmt.Errorf("exceeded max retries: %d", maxRetries)
}

func main() {
	logger := setupLogger()

	tracerProvider, shutdown, err := setupTracer()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to setup tracer")
	}
	defer shutdown()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = retryWithLoggingAndTracing(
		ctx,
		simulateServiceCall,
		maxRetries: 3,
		backoff:    time.Second,
)
	if err != nil {
		logger.Err(err).Msg("final error after retries")
	}
}