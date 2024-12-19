package main

import (
	"fmt"
	"net/url"
)

func getLocalizationFromQuery(query string, defaultLocalization string) string {
	// Parse the query string
	parsedQuery, _ := url.ParseQuery(query)

	// Get the localization parameter
	localization := parsedQuery.Get("lang")

	// Check if the localization parameter is present
	if localization == "" {
		// If not present, return the default localization
		return defaultLocalization
	}

	// If present, return the requested localization
	return localization
}

func main() {
	// Example URLs with query parameters
	urls := []string{
		"http://example.com?lang=en",
		"http://example.com?lang=es",
		"http://example.com",
	}

	// Default localization
	defaultLocalization := "en"

	// Iterate over the URLs and get the localization
	for _, url := range urls {
		localization := getLocalizationFromQuery(url, defaultLocalization)
		fmt.Printf("Localization for URL '%s': %s\n", url, localization)
	}
}
