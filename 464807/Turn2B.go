package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/spf13/viper"
)

func getLocalizationFromQuery(query string) string {
	// Parse the query string
	parsedQuery, _ := url.ParseQuery(query)

	// Get the localization parameter
	localization := parsedQuery.Get("lang")

	// Check if the localization parameter is present
	if localization == "" {
		// If not present, retrieve the default localization from the configuration file
		localization = viper.GetString("default_localization")

		// Handle the case when the default localization is missing in the configuration
		if localization == "" {
			fmt.Println("Error: Default localization not found in configuration.")
			return ""
		}
	}

	return localization
}

func loadConfiguration(configPath string) error {
	// Set the configuration file path
	viper.SetConfigFile(configPath)

	// Read the configuration file
	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("Error reading configuration file: %w", err)
	}

	// Validate that the default localization key exists
	if !viper.IsSet("default_localization") {
		return fmt.Errorf("Missing required key: default_localization in configuration file")
	}

	return nil
}

func main() {
	// Configuration file path (you can change this to a relative path or a command-line flag)
	configPath := "Turn2B.go.yaml"

	// Load the configuration
	err := loadConfiguration(configPath)
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// Example URLs with query parameters
	urls := []string{
		"http://example.com?lang=en",
		"http://example.com?lang=es",
		"http://example.com",
	}

	// Iterate over the URLs and get the localization
	for _, url := range urls {
		localization := getLocalizationFromQuery(url)
		fmt.Printf("Localization for URL '%s': %s\n", url, localization)
	}
}
