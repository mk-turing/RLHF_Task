package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
)

type RateLimitClaim struct {
	jwt.StandardClaims
	Role      string `json:"role"`
	RateLimit int    `json:"rate_limit"`
}

var (
	secretKey   = []byte("your_secret_key") // Should be securely stored
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Define role rate limits
	roleRateLimits = map[string]int{
		"admin":    100, // 100 requests per minute
		"premium":  50,  // 50 requests per minute
		"standard": 20,  // 20 requests per minute
	}
)

func checkRateLimit(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	token, err := jwt.ParseWithClaims(tokenString[7:], &RateLimitClaim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if !token.Valid {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	claims := token.Claims.(*RateLimitClaim)
	userID := claims.Subject
	role := claims.Role

	// Get the rate limit for the user's role
	rateLimit, ok := roleRateLimits[role]
	if !ok {
		http.Error(w, fmt.Sprintf("unknown role: %s", role), http.StatusUnauthorized)
		return
	}

	// Check rate limit in Redis using pipelining
	key := fmt.Sprintf("rate_limit:%s:%s", userID, role)
	pipe := redisClient.TxPipeline()

	incrCmd := pipe.Incr(context.Background(), key)

	pipe.Expire(context.Background(), key, time.Duration(rateLimit*1000*60)) // Rate limit is per minute
	_, err = pipe.Exec(context.Background())
	if err != nil {
		log.Printf("error executing pipeline: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if incrCmd.Err() != nil {
		log.Printf("error getting current count: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if incrCmd.Val() > int64(rateLimit) {
		http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
		return
	}

	// If the rate limit is not exceeded, proceed to the next handler
	next(w, r)
}

func main() {
	http.HandleFunc("/protected", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Protected Resource Accessed!"))
	})

	http.Handle("/", http.StripPrefix("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkRateLimit(w, r, http.DefaultServeMux.ServeHTTP)
	})))

	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("failed to listen and serve: %v", err)
	}
}
