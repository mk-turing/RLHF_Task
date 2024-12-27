package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
)

func sanitizeQueryParameter(queryParam string) string {
	// Pattern to match any non-alphanumeric, non-underscore, and non-dot character
	invalidCharPattern := regexp.MustCompile("[^a-zA-Z0-9_.]+")

	// Replace invalid characters with underscores
	sanitizedParam := invalidCharPattern.ReplaceAllString(queryParam, "_")

	return sanitizedParam
}

// SanitizeQueryParameters middleware sanitizes URL query parameters.
func SanitizeQueryParameters(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the URL query parameters
		query := r.URL.Query()

		// Iterate over each query parameter and sanitize its value
		for paramName, paramValues := range query {
			sanitizedValues := make([]string, 0, len(paramValues))
			for _, paramValue := range paramValues {
				sanitizedValue := sanitizeQueryParameter(paramValue)
				sanitizedValues = append(sanitizedValues, sanitizedValue)
			}

			// Replace the unsanitized values with the sanitized ones
			query[paramName] = sanitizedValues
		}

		// Update the request URL with the sanitized query parameters
		r.URL.RawQuery = query.Encode()

		// Call the next handler in the middleware chain
		next.ServeHTTP(w, r)
	})
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Example: Retrieve a query parameter named "user"
	user := r.URL.Query().Get("user")

	fmt.Fprintf(w, "Hello, %s!", user)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	// Use the SanitizeQueryParameters middleware before the main handler
	sanitizedMux := SanitizeQueryParameters(mux)

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", sanitizedMux))
}
