package main

import (
	"log"
	"net/http"
)

// LogQueryParameters is a generic middleware that logs URL query parameters
func LogQueryParameters(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get query parameters
		params := r.URL.Query()

		// Extract and log each parameter with its type
		log.Printf("Query Parameters:")
		for key, values := range params {
			for _, value := range values {
				// Log the parameter and its value
				log.Printf("Parameter '%s' = %v (type %T)", key, value)
			}
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})

	// Apply the middleware
	http.ListenAndServe(":8080", LogQueryParameters(http.DefaultServeMux))
}
