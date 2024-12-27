package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// Configuration struct to hold filtering rules for each parameter
type Configuration struct {
	FilteringRules map[string]string `json:"filtering_rules"`
}

func (c Configuration) sanitizeParameter(key, value string) string {
	if rule, ok := c.FilteringRules[key]; ok {
		// Apply the specific rule if defined
		return sanitizeWithRule(value, rule)
	} else {
		// Apply the default sanitization rule if not defined
		return sanitize(value)
	}
}

// Example sanitization rules
const (
	defaultRule = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-@.%=:?"
	sqlRule     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_%=:?"
	emailRule   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_.-@"
)

// Enhanced sanitization methods with custom rules
func sanitize(input string) string {
	return sanitizeWithRule(input, defaultRule)
}

func sanitizeWithRule(input, rule string) string {
	// Apply the custom rule for sanitization
	sanitized := strings.Map(func(r rune) rune {
		if strings.ContainsRune(rule, r) {
			return r
		}
		return -1
	}, input)
	return string(sanitized)
}

func sanitizeQueryParameters(next http.Handler, config Configuration) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loggingHandler(w, r, func() {
			q := r.URL.Query()
			for key, values := range q {
				sanitizedValues := []string{}
				for _, value := range values {
					sanitizedValue := config.sanitizeParameter(key, value)
					sanitizedValues = append(sanitizedValues, sanitizedValue)
				}
				q.Set(key, strings.Join(sanitizedValues, ","))
			}
			r.URL.RawQuery = q.Encode()
		})
		next.ServeHTTP(w, r)
	})
}

func loggingHandler(w http.ResponseWriter, r *http.Request, h func()) {
	// Pretend to log and identify potential injection attempts
	h()
}

func main() {
	// Define custom filtering rules for certain parameters
	configuration := Configuration{
		FilteringRules: map[string]string{
			"email":   emailRule,
			"message": sqlRule, // Example of using different rules for specific parameters
		},
	}

	r := mux.NewRouter()
	r.HandleFunc("/example", exampleHandler).Methods("GET")
	// Pass the configuration to the middleware
	r.Use(sanitizeQueryParameters(http.DefaultServeMux, configuration))

	log.Fatal(http.ListenAndServe(":8080", r))
}

func exampleHandler(w http.ResponseWriter, r *http.Request) {
	// Access sanitized query parameters
	q := r.URL.Query()
	email := q.Get("email")
	message := q.Get("message")
	// Respond with the sanitized parameters