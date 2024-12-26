package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Product represents a product in the search system
type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
	Date     string  `json:"date"`
}

// SearchProducts handles the search request
func SearchProducts(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	r.ParseForm()

	// Extract search terms
	searchTerm := strings.Join(r.Form["search"], " ")

	// Extract filters
	filters := make(map[string]string)
	for key, values := range r.Form {
		if len(values) > 0 {
			filters[key] = values[0]
		}
	}

	// Validate and process filters dynamically
	validFilters, err := validateFilters(filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Fetch products from the in-memory list (replace with actual database query)
	products := searchInMemoryProducts(searchTerm, validFilters)

	// Respond with the search results
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// validateFilters validates each filter and returns a map of valid filters
func validateFilters(filters map[string]string) (map[string]interface{}, error) {
	validFilters := make(map[string]interface{})
	for key, value := range filters {
		switch key {
		case "category":
			validFilters[key] = value
		case "date_range":
			dateRange, err := validateDateRange(value)
			if err != nil {
				return nil, err
			}
			validFilters[key] = dateRange
		case "price":
			price, err := validatePrice(value)
			if err != nil {
				return nil, err
			}
			validFilters[key] = price
		default:
			return nil, fmt.Errorf("invalid filter parameter: %s", key)
		}
	}
	return validFilters, nil
}

// validateDateRange validates a date range filter in the format "YYYY-MM-DD,YYYY-MM-DD"
func validateDateRange(dateRange string) ([]time.Time, error) {
	dates := strings.Split(dateRange, ",")
	if len(dates) != 2 {
		return nil, fmt.Errorf("invalid date range format: %s", dateRange)
	}

	var validDates []time.Time
	for _, dateStr := range dates {
		t, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, fmt.Errorf("invalid date format: %s", dateStr)
		}
		validDates = append(validDates, t)
	}

	return validDates, nil
}

// validatePrice validates a price filter in the format "<number>" or "<number>-<number>"
func validatePrice(priceStr string) ([]float64, error) {
	prices := strings.Split(priceStr, "-")
	if len(prices) == 1 {
		// Validate a single price
		price, err := parsePrice(prices[0])
		if err != nil {
			return nil, err
		}
		return []float64{price}, nil
	} else if len(prices) == 2 {
		// Validate a price range
		minPrice, err := parsePrice(prices[0])
		if err != nil {
			return nil, err
		}
		maxPrice, err := parsePrice(prices[1])
		if err != nil {
			return nil, err
		}