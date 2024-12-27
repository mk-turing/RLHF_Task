package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Config struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

var configs = map[string]Config{
	"api-key-1": Config{
		AllowedOrigins: []string{"https://example1.com", "https://sub.example1.com"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type", "Authorization", "X-Custom-Header"},
	},
	"api-key-2": Config{
		AllowedOrigins: []string{"https://example2.com"},
		AllowedMethods: []string{"GET", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type"},
	},
}

func queryParameterInspector(w http.ResponseWriter, r *http.Request) {
	// Parse the URL
	parsedURL, err := url.Parse(r.URL.String())
	if err != nil {
		if err == url.ErrInvalidURL {
			http.Error(w, "Malformed URL. Please check your URL format.", http.StatusBadRequest)
		} else {
			http.Error(w, fmt.Sprintf("Unexpected error parsing URL: %v", err), http.StatusInternalServerError)
		}
		return
	}

	// Parse query parameters from both URL and form
	query := parsedURL.Query()
	query.Add("form-param", r.FormValue("form-param"))

	// Validate query parameters (example: ensuring required fields exist)
	requiredFields := map[string]bool{"param1": true, "form-param": true}
	for field := range requiredFields {
		if len(query[field]) == 0 {
			http.Error(w, fmt.Sprintf("Missing required parameter: %s", field), http.StatusBadRequest)
			return
		}
	}

	// Set CORS headers based on the request headers
	apiKey := r.Header.Get("X-Api-Key")
	config, ok := configs[apiKey]
	if !ok {
		http.Error(w, "Invalid API key. Please provide a valid API key.", http.StatusUnauthorized)
		return
	}

	origin := r.Header.Get("Origin")
	if origin == "" {
		http.Error(w, "No Origin header provided. Please include an Origin header in your request.", http.StatusForbidden)
		return
	}

	if !isAllowedOrigin(config.AllowedOrigins, origin) {
		http.Error(w, fmt.Sprintf("Origin %s not allowed. Please check your domain.", origin), http.StatusForbidden)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ", "))
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ", "))
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Process the request
	// ...
	fmt.Fprintf(w, "Query Parameters Validated: %s", query)
}

func isAllowedOrigin(allowedOrigins []string, origin string) bool {
	for _, allowed := range allowedOrigins {
		if allowed == origin {
			return true
		}
	}
	return false
}

func main() {
	http.HandleFunc("/query", queryParameterInspector)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
