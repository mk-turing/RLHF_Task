package main

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"strconv"
)

var (
	client *redis.Client
	tiers  map[string]int

	requestsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"endpoint", "code"},
	)

	throttlingCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_throttling_total",
			Help: "Total number of HTTP requests throttled",
		},
		[]string{"endpoint", "tier"},
	)
)

func init() {
	prometheus.MustRegister(requestsCounter)
	prometheus.MustRegister(throttlingCounter)

	// Initialize Redis client and tier map as before
	// ...
}

// Logging function to log throttling events
func logThrottlingEvent(user, endpoint, tier, message string) {
	log.Printf("Throttling event: user=%s, endpoint=%s, tier=%s, message=%s", user, endpoint, tier, message)
}

func endpointRateLimit(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ... (Rate-limiting logic as before)

		if count > tierLimit {
			throttlingCounter.WithLabelValues(endpoint, tier).Inc()
			logThrottlingEvent(userID, endpoint, tier, fmt.Sprintf("Rate limit exceeded: %d requests to %s in %s", tierLimit, endpoint, window.String()))

			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		// Run the actual handler and log the request
		handler.ServeHTTP(w, r)
		requestsCounter.WithLabelValues(endpoint, strconv.Itoa(w.WriteHeader(0))).Inc()
	})
}

func main() {
	http.HandleFunc("/", endpointRateLimit(http.HandlerFunc(handler)).ServeHTTP)
	http.Handle("/metrics", promhttp.Handler())

	log.Printf("Server listening on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
