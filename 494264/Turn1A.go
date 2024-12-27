package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// AllowedOrigins represents the list of allowed CORS origins
var AllowedOrigins = []string{"https://example.com", "https://sub.example.com"}

// AllowedMethods represents the list of allowed CORS HTTP methods
var AllowedMethods = []string{"GET", "POST", "PUT", "DELETE"}

// AllowedHeaders represents the list of allowed CORS custom headers
var AllowedHeaders = []string{"Content-Type", "Authorization"}

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

	// Set CORS headers
	origin := r.Header.Get("Origin")
	if origin == "" {
		http.Error(w, "No Origin header provided", http.StatusForbidden)
		return
	}

	if isAllowedOrigin(origin) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(AllowedMethods, ", "))
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(AllowedHeaders, ", "))
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
	} else {
		http.Error(w, fmt.Sprintf("Origin %s not allowed", origin), http.StatusForbidden)
		return
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

func main() {
	http.HandleFunc("/query", queryParameterInspector)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
