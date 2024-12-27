package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func inspectQueryParams(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		allowedOrigins := []string{"https://example.com", "http://localhost:8080"}

		// Check if the origin is allowed
		for _, allowedOrigin := range allowedOrigins {
			if strings.HasPrefix(origin, allowedOrigin) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		// Preflight requests should not have body
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Max-Age", "86400") // Cache preflight results for an hour
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Parse query parameters and inspect them here
		// E.g., check for required parameters, data types, etc.
		err := r.ParseForm()
		if err != nil {
			log.Println("Error parsing query parameters:", err)
			http.Error(w, "Invalid query parameters", http.StatusBadRequest)
			return
		}

		// Pass the request to the next handler if everything looks good
		next.ServeHTTP(w, r)
	})
}

func main() {
	http.Handle("/", inspectQueryParams(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, world! Query parameters: %v", r.URL.Query())
	})))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting HTTP server:", err)
	}
}
