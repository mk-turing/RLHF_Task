package main

import (
	"context"
	"fmt"
	"net/http"
)

var (
	// supportedLocalizations is a map of supported localizations and their corresponding default values
	supportedLocalizations = map[string]string{
		"en": "en", // English
		"es": "es", // Spanish
		"fr": "fr", // French
		"de": "de", // German
	}
	defaultLocalization = "en"
)

func localizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse the URL query parameters
		query := r.URL.Query()
		localizations := query["locale"]

		// Initialize the localization to the default value
		localization := defaultLocalization

		// If localizations are provided in the query, prioritize them
		for _, l := range localizations {
			// Check if the provided localization is supported
			if supported, ok := supportedLocalizations[l]; ok {
				localization = supported
				break
			}
		}

		// Set the localization context for the request
		r = r.WithContext(context.WithValue(r.Context(), "localization", localization))

		// Call the next handler in the middleware chain
		next.ServeHTTP(w, r)
	})
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the localization context from the request
		localization := r.Context().Value("localization").(string)

		fmt.Fprintf(w, "Hello, World! in %s\n", localization)
	})

	// Use the localizationMiddleware for all routes
	http.ListenAndServe(":8080", localizationMiddleware(http.DefaultServeMux))
}
