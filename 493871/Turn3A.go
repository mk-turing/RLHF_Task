package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

type Config struct {
	ParamRules map[string]string `json:"param_rules"`
}

func loadConfig() (*Config, error) {
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = "turn3A_config.json"
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

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/example", exampleHandler).Methods("GET")

	// Apply the sanitization middleware to all routes
	r.Use(sanitizeQueryParameters)

	http.ListenAndServe(":8080", r)
}

func exampleHandler(w http.ResponseWriter, r *http.Request) {
	// Access sanitized query parameters
	q := r.URL.Query()
	param := q.Get("param")

	// Respond with the sanitized parameter
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Sanitized Parameter: " + param))
}
