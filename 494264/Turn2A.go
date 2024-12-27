package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// AllowedOrigins represents a fallback list of allowed CORS origins
var AllowedOrigins = []string{"https://example.com"}

// AllowedMethods represents the list of allowed CORS HTTP methods
var AllowedMethods = []string{"GET", "POST", "PUT", "DELETE"}

// AllowedHeaders represents the list of allowed CORS custom headers
var AllowedHeaders = []string{"Content-Type", "Authorization"}

// GlobalConfig contains dynamic CORS configuration configurations
var GlobalConfig = map[string]bool{}

func queryParameterInspector(w http.ResponseWriter, r *http.Request) {
	// Parse the URL
	parsedURL, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	// Parse query parameters
	query := parsedURL.Query()

	// Validate query parameters (example: ensuring required fields exist)
	requiredFields := map[string]bool{"param1": true, "param2": true}
	for field := range requiredFields {
		if len(query[field]) == 0 {
			http.Error(w, fmt.Sprintf("Missing required parameter: %s", field), http.StatusBadRequest)
			return
		}
	}

	// Set CORS headers based on either fallback or custom request origin
	origin := r.Header.Get("Origin")
	if origin == "" {
		origin = AllowedOrigins[0] // Default to the first fallback origin
	}

	// Check if the origin has a dynamic configuration override
	if configOverride, ok := GlobalConfig[origin]; ok {
		if configOverride {
			// The origin is allowed by custom configuration
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(AllowedMethods, ", "))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(AllowedHeaders, ", "))
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
		} else {
			// The origin is explicitly disallowed by custom configuration
			http.Error(w, fmt.Sprintf("Origin %s not allowed", origin), http.StatusForbidden)
			return
		}
	} else {
		// Check fallback AllowedOrigins
		if isAllowedOrigin(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(AllowedMethods, ", "))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(AllowedHeaders, ", "))
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
		} else {
			// Not in fallback list
			http.Error(w, fmt.Sprintf("Origin %s not allowed", origin), http.StatusForbidden)
			return
		}
	}

	// Process the request
	// ...
	fmt.Fprintf(w, "Query Parameters Validated: %s", query)
}

func isAllowedOrigin(origin string) bool {
	for _, allowed := range AllowedOrigins {
		if allowed == origin {
			return true
		}
	}
	return false
}

// Example of setting dynamic configuration based on another header (e.g., API key or authentication)
func SetDynamicCORSConfiguration(apiKey string, isAllowed bool) {
	GlobalConfig[apiKey] = isAllowed
	fmt.Printf("CORS configuration set for API key '%s': %v\n", apiKey, isAllowed)
}

func main() {
	// Example of initializing dynamic configuration
	SetDynamicCORSConfiguration("some-api-key-123", true)
	SetDynamicCORSConfiguration("another-api-key-456", false)

	http.HandleFunc("/query", queryParameterInspector)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
