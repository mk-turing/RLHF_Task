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

func init() {
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0, // use default DB
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic("could not connect to Redis: " + err.Error())
	}
}

func parseQueryParameters(r *http.Request) (limit int, window time.Duration, err error) {
	limitStr := r.URL.Query().Get("limit")
	windowStr := r.URL.Query().Get("window")

	if limitStr == "" || windowStr == "" {
		return 0, 0, fmt.Errorf("rate limit configuration missing in query parameters")
	}

	limit, err = strconv.Atoi(limitStr)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid limit parameter: %v", err)
	}

	window, err = time.ParseDuration(windowStr)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid window parameter: %v", err)
	}

	return
}

func rateLimit(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limit, window, err := parseQueryParameters(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		key := fmt.Sprintf("rate-limit:%s", r.RemoteAddr)

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

		if count > limit {
			// Exceeded limit
			// Set an expiry on the key for the window duration
			err = client.Expire(context.Background(), key, window).Err()
			if err != nil {
				log.Printf("error setting Redis expiry: %v", err)
			}
			http.Error(w, fmt.Sprintf("Rate limit exceeded: %d requests in %s", limit, window.String()), http.StatusTooManyRequests)
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
	http.HandleFunc("/", rateLimit(http.HandlerFunc(handler)).ServeHTTP)
	log.Printf("Server listening on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
