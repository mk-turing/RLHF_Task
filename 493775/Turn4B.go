package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
	Date     string  `json:"date"`
}

type SearchResponse struct {
	Products []Product `json:"products"`
	Total    int       `json:"total"`
	Page     int       `json:"page"`
	Pages    int       `json:"pages"`
	Size     int       `json:"size"`
}

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

	// Validate filters
	validFilters := validateFilters(filters)

	// Parse page and size parameters
	page := 1
	size := 10
	if value, ok := r.Form["page"]; ok {
		if p, err := strconv.Atoi(value[0]); err == nil && p > 0 {
			page = p
		}
	}
	if value, ok := r.Form["size"]; ok {
		if s, err := strconv.Atoi(value[0]); err == nil && s > 0 {
			size = s
		}
	}

	// Fetch products from the in-memory list (replace with actual database query)
	allProducts := searchInMemoryProducts(searchTerm, validFilters)

	// Paginate the results
	total := len(allProducts)
	pages := (total + size - 1) / size
	if page > pages {
		page = pages
	}
	start := (page - 1) * size
	end := min(start+size, total)
	products := allProducts[start:end]

	// Prepare the response with pagination information
	response := SearchResponse{
		Products: products,
		Total:    total,
		Page:     page,
		Pages:    pages,
		Size:     size,
	}

	// Respond with the search results
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// rest of the code remains the same

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

// sortProductsByMultipleCriteria sorts products by multiple criteria
func sortProductsByMultipleCriteria(products []Product, sortFields []string, sortOrders []string) []Product {
	sort.Slice(products, func(i, j int) bool {
		for k, field := range sortFields {
			order := sortOrders[k]
			switch field {
			case "name":
				if order == "asc" {
					if products[i].Name < products[j].Name {
						return true
					} else if products[i].Name > products[j].Name {
						return false
					}
				} else if order == "desc" {
					if products[i].Name > products[j].Name {
						return true
					} else if products[i].Name < products[j].Name {
						return false
					}
				}
			case "price":
				if order == "asc" {
					if products[i].Price < products[j].Price {
						return true
					} else if products[i].Price > products[j].Price {
						return false
					}
				} else if order == "desc" {
					if products[i].Price > products[j].Price {
						return true
					} else if products[i].Price < products[j].Price {
						return false
					}
				}
			case "date":
				if order == "asc" {
					if products[i].Date < products[j].Date {
						return true
					} else if products[i].Date > products[j].Date {
						return false
					}
				} else if order == "desc" {
					if products[i].Date > products[j].Date {
						return true
					} else if products[i].Date < products[j].Date {
						return false
					}
				}
			}
		}
		return false
	})
	return products
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
