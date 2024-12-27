package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

func sanitizeQueryParameters(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse the query string
		q := r.URL.Query()

		// Sanitize each query parameter
		for key, values := range q {
			sanitizedValues := []string{}
			for _, value := range values {
				sanitizedValue := sanitize(value)
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

func sanitize(input string) string {
	// Basic sanitization: remove non-alphanumeric characters, spaces, and special characters
	// This is a simplified example and might need further customization based on the application's needs
	sanitized := strings.Map(func(r rune) rune {
		if strings.ContainsRune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_.- ", r) {
			return r
		}
		return -1
	}, input)

	return string(sanitized)
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
