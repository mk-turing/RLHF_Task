package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
	"strconv"
	"time"
)

var client *redis.Client
var tiers map[string]int

func init() {
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0, // use default DB
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic("could not connect to Redis: " + err.Error())
	}

	tiers = map[string]int{
		"basic":   10,  // Basic user: 10 requests per minute
		"premium": 50,  // Premium user: 50 requests per minute
		"elite":   100, // Elite user: 100 requests per minute
	}
}

func parseQueryParameters(r *http.Request) (tier string, limit int, window time.Duration, err error) {
	tier = r.URL.Query().Get("tier")
	limitStr := r.URL.Query().Get("limit")
	windowStr := r.URL.Query().Get("window")

	if tier == "" || limitStr == "" || windowStr == "" {
		return "", 0, 0, fmt.Errorf("required query parameters missing")
	}

	limit, err = strconv.Atoi(limitStr)
	if err != nil {
		return "", 0, 0, fmt.Errorf("invalid limit parameter: %v", err)
	}

	window, err = time.ParseDuration(windowStr)
	if err != nil {
		return "", 0, 0, fmt.Errorf("invalid window parameter: %v", err)
	}

	return
}

func getUserID(r *http.Request) string {
	// For demonstration purposes, we're using the remote IP as the user ID
	return r.RemoteAddr
}

func endpointRateLimit(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tier, limit, window, err := parseQueryParameters(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userID := getUserID(r)
		endpoint := r.URL.Path
		key := fmt.Sprintf("rate-limit:%s:%s:%s", userID, tier, endpoint)

		// Determine valid rate limit based on tier
		var tierLimit int
		if tl, ok := tiers[tier]; ok {
			tierLimit = tl
		} else {
			tierLimit = limit
			log.Printf("Invalid tier specified, defaulting to user-defined limit: %d", limit)
		}

		// Use Redis to track requests
		_, err = client.IncrBy(context.Background(), key, 1).Result()
		if err != nil {
			log.Printf("error incrementing Redis key: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		count, err := client.Get(context.Background(), key).Int()
		if err != nil {
			log.Printf("error getting Redis key: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if count > tierLimit {
			// Exceeded limit
			// Set an expiry on the key for the window duration
			err = client.Expire(context.Background(), key, window).Err()
			if err != nil {
				log.Printf("error setting Redis expiry: %v", err)
			}
			http.Error(w, fmt.Sprintf("Rate limit for %s tier exceeded: %d requests to %s in %s", tier, tierLimit, endpoint, window.String()), http.StatusTooManyRequests)
			return
		}

		// Run the actual handler
		handler.ServeHTTP(w, r)
	})
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, Rate Limited World!")
}

func main() {
	http.HandleFunc("/", endpointRateLimit(http.HandlerFunc(handler)).ServeHTTP)
	log.Printf("Server listening on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
