package main

import (
	"context"
	"fmt"
	"github.com/ulule/limiter/v3"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// Define your custom claim struct to hold the rate limit information
type RateLimitClaims struct {
	jwt.StandardClaims
	RateLimit int `json:"rate_limit"`
}

func main() {
	// Create a new limiter store with a memory backend
	store := memory.NewStore()

	// Create a new rate limiter for the specified key and rate
	rate := limiter.Rate{Limit: 5, Period: time.Second}
	key := "user-1" // Replace this with the actual user ID or any other key
	rateLimiter, err := store.Get(key)
	if err != nil {
		log.Fatal(err)
	}
	rateLimiter.SetLimit(rate)

	// Create a new JWT signing key
	signingKey := []byte("secret")

	// Create a middleware function to validate rate limit using JWT claims
	rateLimitMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract the JWT token from the request header
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				http.Error(w, "JWT token is required", http.StatusUnauthorized)
				return
			}

			// Parse the JWT token and retrieve the rate limit claim
			token, err := jwt.ParseWithClaims(tokenString, &RateLimitClaims{}, func(token *jwt.Token) (interface{}, error) {
				return signingKey, nil
			})
			if err != nil {
				http.Error(w, "Invalid JWT token", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(*RateLimitClaims)
			if !ok || !token.Valid {
				http.Error(w, "Invalid JWT claims", http.StatusUnauthorized)
				return
			}

			// Create a new context with the rate limit claim
			ctx := context.WithValue(r.Context(), "rateLimit", claims.RateLimit)

			// Apply the rate limiter middleware
			stdlib.RateLimit(rateLimiter, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r.WithContext(ctx))
			})).ServeHTTP(w, r)
		})
	}

	// Sample route handler to demonstrate rate limiting
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rateLimit := r.Context().Value("rateLimit").(int)
		fmt.Fprintf(w, "Hello! Your rate limit is: %d requests per second.\n", rateLimit)
	})

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", rateLimitMiddleware(http.DefaultServeMux)))
}
