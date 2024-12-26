package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Product represents a product in our search system
type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
	Date     string  `json:"date"`
}

func validateFilters(filters map[string][]string) map[string][]string {
	validFilters := make(map[string][]string)
	for key, values := range filters {
		switch key {
		case "category":
			// Validate that categories are strings
			validCategories := []string{}
			for _, value := range values {
				if _, ok := validFilters[key]; ok {
					continue
				}
				if isValidString(value) {
					validCategories = append(validCategories, value)
				}
			}
			if len(validCategories) > 0 {
				validFilters[key] = validCategories
			}
		case "date":
			// Validate that dates are in YYYY-MM-DD format
			validDates := []string{}
			for _, value := range values {
				if isValidDate(value) {
					validDates = append(validDates, value)
				}
			}
			if len(validDates) > 0 {
				validFilters[key] = validDates
			}
		case "price":
			// Validate that prices are numeric
			validPrices := []string{}
			for _, value := range values {
				if isValidPrice(value) {
					validPrices = append(validPrices, value)
				}
			}
			if len(validPrices) > 0 {
				validFilters[key] = validPrices
			}
		}
	}
	return validFilters
}

func isValidPrice(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// SearchProducts handles the search request
func SearchProducts(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	r.ParseForm()

	// Extract search terms
	searchTerm := strings.Join(r.Form["search"], " ")

	// Extract filters
	filters := make(map[string][]string)
	for _, value := range r.Form["filter"] {
		parts := strings.SplitN(value, ":", 2)
		if len(parts) == 2 {
			if filters[parts[0]] == nil {
				filters[parts[0]] = make([]string, 0)
			}
			filters[parts[0]] = append(filters[parts[0]], parts[1])
		}
	}

	// Extract sorting options
	sortFields := r.Form["sort"]
	order := r.FormValue("order")

	// Validate filters and sorting options
	validFilters := validateFilters(filters)
	validSortFields, validOrder := validateSortingOptions(sortFields, order)

	// Fetch products from the in-memory list (replace with actual database query)
	products := searchInMemoryProducts(searchTerm, validFilters)

	// Sort products based on the specified criteria
	products = sortProducts(products, validSortFields, validOrder)

	// Respond with the search results
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// searchInMemoryProducts performs a search on an in-memory list of products (remains the same)

// sortProducts performs multi-field sorting
func sortProducts(products []Product, sortFields []string, order string) []Product {
	if len(sortFields) == 0 {
		return products // Return the original list if no sorting is specified
	}

	lessFunctions := make([]func(Product, Product) bool, len(sortFields))

	for i, sortField := range sortFields {
		switch sortField {
		case "price":
			lessFunctions[i] = sortByPrice(order == "asc")
		case "date":
			lessFunctions[i] = sortByDate(order == "asc")
		case "relevance":
			lessFunctions[i] = sortByRelevance(order == "asc")
		default:
			// Handle invalid sort fields (this will never happen as we validate the sorting options earlier)
			panic(fmt.Sprintf("Invalid sort field: %s", sortField))
		}
	}

	for i := len(sortFields) - 1; i > 0; i-- {
		for j := 0; j < i; j++ {
			if lessFunctions[i](products[j], products[j+1]) {
				products[j], products[j+1] = products[j+1], products[j]
			}
		}
	}

	return products
}

// Custom sort functions for different fields
func sortByPrice(asc bool) func(Product, Product) bool {
	return func(p1, p2 Product) bool {
		if asc {
			return p1.Price < p2.Price
		}
		return p1.Price > p2.Price
	}
}

func sortByDate(asc bool) func(Product, Product) bool {
	return func(p1, p2 Product) bool {
		if asc {
			return p1.Date < p2.Date
		}
		return p1.Date > p2.Date
	}
}

func sortByRelevance(asc bool) func(Product, Product) bool {
	// For simplicity, we'll assume relevance is calculated as 1.0 for all products in this example
	return func(p1, p2 Product) bool {
		if asc {
			return p1.Price < p2.Price
		}
		return p1.Price > p2.Price
	}
}

// Validate filters and sorting options (remains the same)

func main() {
	http.HandleFunc("/search", SearchProducts)
	fmt.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
