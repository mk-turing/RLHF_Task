package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
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

type UserTier struct {
	Name   string
	Limit  int
	Window time.Duration
}

var userTiers = map[string]UserTier{
	"free":    {Name: "Free", Limit: 10, Window: time.Minute},
	"basic":   {Name: "Basic", Limit: 50, Window: time.Minute},
	"premium": {Name: "Premium", Limit: 200, Window: time.Minute},
}

func parseQueryParameters(r *http.Request) (tier string, err error) {
	tier = r.URL.Query().Get("tier")
	if tier == "" {
		return "", fmt.Errorf("user tier configuration missing in query parameters")
	}
	_, ok := userTiers[tier]
	if !ok {
		return "", fmt.Errorf("invalid user tier: %s", tier)
	}
	return tier, nil
}

func rateLimit(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tier, err := parseQueryParameters(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userTier := userTiers[tier]

		key := fmt.Sprintf("rate-limit:%s:%s", tier, r.RemoteAddr)

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

		if count > userTier.Limit {
			err = client.Expire(context.Background(), key, userTier.Window).Err()
			if err != nil {
				log.Printf("error setting Redis expiry: %v", err)
			}

			notification := fmt.Sprintf("Rate limit exceeded for tier '%s': %d requests in %s", userTier.Name, userTier.Limit, userTier.Window.String())
			log.Println(notification)

			http.Error(w, notification, http.StatusTooManyRequests)
			return
		}

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
