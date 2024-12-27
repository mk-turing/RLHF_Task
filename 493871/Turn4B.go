package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Config struct {
	ParamRules map[string]string `json:"param_rules"`
	RateLimit  int               `json:"rate_limit"`
	BurstLimit int               `json:"burst_limit"`
}

type limiter struct {
	sync.Mutex
	limiter    *time.Ticker
	rateLimit  int
	burstLimit int
	requests   int
}

func (l *limiter) Allow() bool {
	l.Lock()
	defer l.Unlock()

	if l.limiter == nil {
		l.limiter = time.NewTicker(time.Second)
		go func() {
			for range l.limiter.C {
				l.Lock()
				l.requests = 0
				l.Unlock()
			}
		}()
	}

	if l.requests >= l.burstLimit {
		return false
	}

	l.requests++
	return true
}

func loadConfig() (*Config, error) {
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = "config.json"
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

func sanitizeQueryParameters(next http.Handler) http.Handler {
	config, err := loadConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("Error loading configuration: %v", err))
		return next
	}

	// Initialize the rate limiter
	limiter := &limiter{
		rateLimit:  config.RateLimit,
		burstLimit: config.BurstLimit,
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		// Parse the query string
		q := r.URL.Query()

		// Sanitize each query parameter
		for key, values := range q {
			sanitizedValues := []string{}
			rule := config.ParamRules[key]
			for _, value := range values {
				sanitizedValue := sanitize(value, rule)