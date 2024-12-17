package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"golang.org/x/oauth2"
)

type RateLimitClaim struct {
	jwt.StandardClaims
	Role      string `json:"role"`
	RateLimit int    `json:"rate_limit"`
}

// ConfigurationService struct to fetch rate limits
type ConfigurationService struct {
	client    *http.Client
	configURL string
}

func (c *ConfigurationService) FetchRateLimits() (map[string]int, error) {
	resp, err := c.client.Get(c.configURL)
	if err != nil {
		log.Printf("error fetching rate limits: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var rateLimits map[string]int
	err = json.NewDecoder(resp.Body).Decode(&rateLimits)
	if err != nil {
		log.Printf("error decoding rate limits: %v", err)
		return nil, err
	}

	return rateLimits, nil
}

func newConfigurationService(configURL string, oauthClient *http.Client) *ConfigurationService {
	return &ConfigurationService{
		client:    oauthClient,
		configURL: configURL,
	}
}

var (
	secretKey   = []byte("your_secret_key")
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	// Lua increment script
	luaIncrScript = redis.NewScript(`
		local key = KEYS[1]
		local rateLimit = ARGV[1]
		local expireTime = ARGV[2]

		local count = redis.call("INCR", key)
		if count == 1 then
			redis.call("EXPIRE", key, expireTime)
		end

		if count > rateLimit then
			redis.call("DECR", key)
			return -1
		else
			return count
		end
	`)

	// Role rate limits cache key
	rateLimitsCacheKey = "role_rate_limits"

	// OAuth client
	oauthClient      *http.Client
	configServiceURL = "http://example.com/rate_limits"
)

func updateRoleRateLimitsCache(ctx context.Context, configService *ConfigurationService) {
	for {
		select {
		case <-time.After(time.Minute * 5): // Fetch every 5 minutes
			log.Println("Fetching rate limits from configuration service...")
			newRateLimits, err := configService.FetchRateLimits()
			if err != nil {
				log.Printf("Failed to fetch rate limits: %v", err)
				continue
			}

			// Store new rate limits in Redis
			for role, limit := range newRateLimits {
				if err := redisClient.HSet(ctx, rateLimitsCacheKey, role, limit).Err(); err != nil {
					log.Printf("Error setting new rate limits in cache: %v", err)
				}
			}

			log.Println("Rate limits cache updated successfully.")
		case <-ctx.Done():
			log.Println("Stopping role rate limits cache update loop.")
			return
		}
	}
}

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

	// Fetch the role rate limits from Redis cache
	var fetchedRoleRateLimits map[string]string
	fetchedRoleRateLimits, err = redisClient.HGetAll(context.Background(), rateLimitsCacheKey).Result()
	if err != nil {
		log.Printf("error fetching role rate limits from cache: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	rateLimit, ok := fetchedRoleRateLimits[role]
	if !ok {
		http.Error(w, fmt.Sprintf("unknown role: %s", role), http.StatusUnauthorized)
		return
	}

	// Convert rateLimit to int
	rateLimitInt, err := strconv.Atoi(rateLimit)
	if err != nil {
		log.Printf("Error converting rate limit to int: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Calculate the expiration time for the rate limit key in seconds
	expireTime := int64(rateLimitInt * 60)

	// Use LUA script to increment the rate limit count atomically
	key := fmt.Sprintf("rate_limit:%s:%s", userID, role)
	count, err := luaIncrScript.Run(context.Background(), redisClient, []string{key}, rateLimitInt, expireTime).Result()
	if err != nil {
		log.Printf("error running LUA script: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if count.(int64) == -1 {
		http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
		return
	}

	// If the rate limit is not exceeded, proceed to the next handler
	next(w, r)
}

func main() {
	// Initialize OAuth client (example, replace with your own)
	oauthClient = oauth2.NewClient(context.Background(), nil) // replace with actual OAuth2 token source

	// Initialize config service
	configService := newConfigurationService(configServiceURL, oauthClient)

	// Start the background cache update loop
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go updateRoleRateLimitsCache(ctx, configService)

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
