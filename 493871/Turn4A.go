package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Config struct {
	ParamRules        map[string]string `json:"param_rules"`
	RateLimitRequests int               `json:"rate_limit_requests"`
	RateLimitWindow   time.Duration     `json:"rate_limit_window"`
}

func loadConfig() (*Config, error) {
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = "turn4A_config.json"
	}

	file, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func sanitize(input string, rule string) string {
	if rule == "" {
		// Default rule: remove non-alphanumeric characters, spaces, and special characters
		return strings.Map(func(r rune) rune {
			if strings.ContainsRune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-@.%=:?", r) {
				return r
			}
			return -1
		}, input)
	}

	// Example custom rule: allow additional characters
	return strings.Map(func(r rune) rune {
		if strings.ContainsRune(rule+"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-@.%=:?", r) {
			return r
		}
		return -1
	}, input)
}

type RateLimiter struct {
	requests     int
	requestLimit int
	window       time.Duration
	last         time.Time
}

func NewRateLimiter(requestLimit int, window time.Duration) *RateLimiter {
	return &RateLimiter{requestLimit: requestLimit, window: window, last: time.Now()}
}

func (r *RateLimiter) check(add bool) bool {
	// Check if we have hit the rate limit
	if time.Since(r.last) < r.window {
		if add && r.requestLimit == 0 {
			return true
		}
		r.requests += 1
		return r.requests >= r.requestLimit
	}
	r.requests = 0
	r.last = time.Now()
	if add {
		r.requests++
	}
	return false
}

func sanitizeQueryParameters(config *Config, rateLimiter *RateLimiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if request exceeds rate limit
		if rateLimiter.check(true) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		defer rateLimiter.check(false)

		// Parse the query string
		q := r.URL.Query()

		// Sanitize each query parameter
		for key, values := range q {
			sanitizedValues := []string{}
			rule := config.ParamRules[key]
			for _, value := range values {
				sanitizedValue := sanitize(value, rule)
				sanitizedValues = append(sanitizedValues, sanitizedValue)
			}
			q.Set(key, strings.Join(sanitizedValues, ","))
		}

		// Rebuild the URL with sanitized query parameters
		r.URL.RawQuery = q.Encode()

		// Handle POST request form data
		if r.Method == http.MethodPost {
			r.ParseForm()
			for key, values := range r.PostForm {
				sanitizedValues := []string{}
				rule := config.ParamRules[key]
				for _, value := range values {
					sanitizedValue := sanitize(value, rule)
					sanitizedValues = append(sanitizedValues, sanitizedValue)
				}
				r.PostForm.Set(key, strings.Join(sanitizedValues, ","))
			}
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/example", exampleHandler).Methods("GET", "POST")

	config, err := loadConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("Error loading configuration: %v", err))
	}

	// Create a rate limiter
	rateLimiter := NewRateLimiter(config.RateLimitRequests, config.RateLimitWindow*time.Second)

	// Apply the sanitization middleware to all routes
	r.Use(func(next http.Handler) http.Handler {
		return sanitizeQueryParameters(config, rateLimiter, next)
	})

	// Enable CORS for testing
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST"},
	})

	http.ListenAndServe(":8080", c.Handler(r))
}

func exampleHandler(w http.ResponseWriter, r *http.Request) {
	// Access sanitized query parameters
	q := r.URL.Query()
	param := q.Get("param")

	// Access sanitized form data in POST requests
	if r.Method == http.MethodPost {
		postParam := r.PostForm.Get("postParam")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Sanitized Post Parameter: " + postParam))
		return
	}

	// Respond with the sanitized query parameter
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Sanitized Parameter: " + param))
}
