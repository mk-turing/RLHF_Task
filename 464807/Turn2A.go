package main

import (
	"fmt"
	"net/url"

	"github.com/spf13/viper"
)

func main() {
	// Set up Viper for configuration
	viper.SetConfigFile("Turn2A.yaml") // Specify the configuration file path
	viper.SetConfigType("yaml")        // Specify the configuration file format

	// Load configuration
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Configuration file not found")
		} else {
			fmt.Println("Error reading configuration file:", err)
		}
		return
	}

	// Get the default localization from the configuration
	defaultLocalization := viper.GetString("default_lang")

	if defaultLocalization == "" {
		fmt.Println("Default localization not set in configuration file")
		return
	}

	// Example URLs with query parameters
	urls := []string{
		"http://example.com?lang=en",
		"http://example.com?lang=es",
		"http://example.com",
	}

	// Iterate over the URLs and get the localization
	for _, url := range urls {
		localization := getLocalizationFromQuery(url, defaultLocalization)
		fmt.Printf("Localization for URL '%s': %s\n", url, localization)
	}
}

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
