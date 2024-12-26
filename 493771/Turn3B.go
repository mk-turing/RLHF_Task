package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var client *redis.Client
var tiers map[string]int
var endpointLimits map[string]int

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
		"basic":   10,
		"premium": 50,
		"elite":   100,
	}
	endpointLimits = map[string]int{
		"/":       10, // Default limit for root endpoint
		"/api/v1": 20, // Limit for /api/v1 endpoint
	}
}

func parseQueryParameters(r *http.Request) (tier, endpoint string, limit, window int, err error) {
	tier = r.URL.Query().Get("tier")
	endpoint = r.URL.Query().Get("endpoint")
	limitStr := r.URL.Query().Get("limit")
	windowStr := r.URL.Query().Get("window")

	if tier == "" || endpoint == "" || limitStr == "" || windowStr == "" {
		return "", "", 0, 0, fmt.Errorf("required query parameters missing")
	}

	limit, err = strconv.Atoi(limitStr)
	if err != nil {
		return "", "", 0, 0, fmt.Errorf("invalid limit parameter: %v", err)
	}

	window, err = strconv.Atoi(windowStr)
	if err != nil {
		return "", "", 0, 0, fmt.Errorf("invalid window parameter: %v", err)
	}

	return
}

type rateLimiter struct {
	sync.Mutex
	count   int
	expires time.Time
}

func getRateLimiter(key string) *rateLimiter {
	rl := &rateLimiter{}
	data, err := client.Get(context.Background(), key).Bytes()
	if err == nil {
		err = rl.UnmarshalBinary(data)
		if err != nil {
			log.Printf("error unmarshaling rate limiter data: %v", err)
			return nil // Return nil to handle errors gracefully and fallback to default behavior
		}
	}
	return rl
}

func saveRateLimiter(key string, rl *rateLimiter) {
	data, err := rl.MarshalBinary()
	if err != nil {
		log.Printf("error marshaling rate limiter data: %v", err)
		return
	}
	err = client.Set(context.Background(), key, data, rl.expires.Sub(time.Now())).Err()
	if err != nil {
		log.Printf("error saving rate limiter to Redis: %v", err)
	}
}

func (rl *rateLimiter) Increment(limit int, window time.Duration) bool {
	rl.Lock()
	defer rl.Unlock()

	now := time.Now()
	if rl.count == 0 || rl.expires.Before(now) {
		rl.count = 1
		rl.expires = now.Add(window)
	} else {
		rl.count++
	}

	return rl.count <= limit
}

func tierRateLimit(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {