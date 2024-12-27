package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var (
	logger = logrus.New()

	validationSuccess = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "query_validation_success",
			Help: "Total number of successful query parameter validations.",
		},
		[]string{"parameter", "value"},
	)

	validationErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "query_validation_errors",
			Help: "Total number of query parameter validation errors.",
		},
		[]string{"parameter", "value", "error"},
	)

	rateLimitedRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "query_rate_limited_requests",
			Help: "Total number of rate-limited requests based on query parameters.",
		},
		[]string{"parameter", "value"},
	)
)

func init() {
	prometheus.MustRegister(validationSuccess)
	prometheus.MustRegister(validationErrors)
	prometheus.MustRegister(rateLimitedRequests)
}

type queryValidator struct {
	params      map[string][]validator
	rateLimiter *rateLimiter
}

// Rest of the code remains the same...

func logValidationResult(status string, param, value, error string) {
	fields := logrus.Fields{"param": param, "value": value}
	if status == "success" {
		logger.WithFields(fields).Info("Query parameter validation succeeded.")
		validationSuccess.With(fields).Inc()
	} else {
		fields["error"] = error
		logger.WithFields(fields).Error("Query parameter validation failed.")
		validationErrors.With(fields).Inc()
	}
}

func logRateLimitingAction(param, value string) {
	logger.WithFields(logrus.Fields{"param": param, "value": value}).Info("Query parameter exceeded rate limit, request rejected.")
	rateLimitedRequests.With(logrus.Fields{"param": param, "value": value}).Inc()
}

func main() {
	// Logging configuration (using Logrus for simplicity)
	logger.Out = logrus.StandardLogger().Out
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.DebugLevel) // Example setting log level to Debug

	// Rest of the code remains the same...

	http.Handle("/api/data", v.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Data endpoint reached")
	})))

	// Metrics endpoint for Prometheus scraping
	http.Handle("/metrics", promhttp.Handler())

	log.Fatal(http.ListenAndServe(":8080", nil))
}
