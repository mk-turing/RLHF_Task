package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
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
	sort := r.FormValue("sort")
	order := r.FormValue("order")

	// Validate filters
	validFilters := validateFilters(filters)

	// Fetch products from the in-memory list (replace with actual database query)
	products := searchInMemoryProducts(searchTerm, validFilters)

	// Sort products based on the specified criteria
	if sort == "price" {
		products = sortProductsByPrice(products, order)
	}

	// Respond with the search results
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// searchInMemoryProducts performs a search on an in-memory list of products
func searchInMemoryProducts(searchTerm string, filters map[string][]string) []Product {
	products := []Product{
		{ID: 1, Name: "Laptop", Category: "electronics", Price: 999.99, Date: "2023-01-15"},
		{ID: 2, Name: "Smartphone", Category: "electronics", Price: 499.99, Date: "2023-02-20"},
		{ID: 3, Name: "Headphones", Category: "electronics", Price: 199.99, Date: "2023-03-10"},
		{ID: 4, Name: "Refrigerator", Category: "appliances", Price: 1499.99, Date: "2023-04-05"},
		{ID: 5, Name: "Washing Machine", Category: "appliances", Price: 899.99, Date: "2023-05-18"},
	}

	// Filter products based on criteria
	var filteredProducts []Product
	for _, product := range products {
		match := true
		for filterKey, filterValues := range filters {
			if filterKey == "category" && !contains(filterValues, product.Category) {
				match = false
				break
			} else if filterKey == "date" {
				for _, dateRange := range filterValues {
					if !isInDateRange(product.Date, dateRange) {
						match = false
						break
					}
				}
				if !match {
					break
				}
			} else if filterKey == "price" {
				for _, priceRange := range filterValues {
					if !isInPriceRange(product.Price, priceRange) {
						match = false
						break
					}
				}
				if !match {
					break
				}
			}
		}
		if match {
			filteredProducts = append(filteredProducts, product)
		}
	}

	// Search for products containing the search term
	var searchResults []Product
	for _, product := range filteredProducts {
		if strings.Contains(strings.ToLower(product.Name), strings.ToLower(searchTerm)) {
			searchResults = append(searchResults, product)
		}
	}

	return searchResults
}

// sortProductsByPrice sorts products based on price
func sortProductsByPrice(products []Product, order string) []Product {
	sortedProducts := make([]Product, len(products))
	copy(sortedProducts, products)

	switch order {
	case "asc":
		for i := 0; i < len(sortedProducts)-1; i++ {
			for j := 0; j < len(sortedProducts)-i-1; j++ {
				if sortedProducts[j].Price > sortedProducts[j+1].Price {
					sortedProducts[j], sortedProducts[j+1] = sortedProducts[j+1], sortedProducts[j]
				}
			}
		}
	case "desc":
		for i := 0; i < len(sortedProducts)-1; i++ {
			for j := 0; j < len(sortedProducts)-i-1; j++ {
				if sortedProducts[j].Price < sortedProducts[j+1].Price {
					sortedProducts[j], sortedProducts[j+1] = sortedProducts[j+1], sortedProducts[j]
				}
			}
		}
	}

	return sortedProducts
}

// validateFilters validates the filters
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

// isValidString checks if the string is non-empty
func isValidString(s string) bool {
	return s != ""
}

// isValidDate checks if the date is in YYYY-MM-DD format
func isValidDate(s string) bool {
	datePattern := `\d{4}-\d{2}-\d{2}`
	return regexp.MustCompile(datePattern).MatchString(s)
}

// isValidPrice checks if the price is a valid float64
func isValidPrice(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// contains checks if a slice contains a string
func contains(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

// isInDateRange checks if a date falls within a range
func isInDateRange(date, dateRange string) bool {
	parts := strings.Split(dateRange, "/")
	if len(parts) != 2 {
		return false // Invalid range format
	}
	start, end := parts[0], parts[1]

	return date >= start && date <= end
}

// isInPriceRange checks if a price falls within a range
func isInPriceRange(price float64, priceRange string) bool {
	parts := strings.Split(priceRange, "/")
	if len(parts) != 2 {
		return false // Invalid range format
	}
	start, _ := strconv.ParseFloat(parts[0], 64)
	end, _ := strconv.ParseFloat(parts[1], 64)

	return price >= start && price <= end
}

func main() {
	http.HandleFunc("/search", SearchProducts)
	fmt.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
