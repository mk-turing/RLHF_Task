package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	filters := make(map[string]string)
	for _, value := range r.Form["filter"] {
		parts := strings.SplitN(value, ":", 2)
		if len(parts) == 2 {
			filters[parts[0]] = parts[1]
		}
	}

	// Extract sorting options
	sort := r.FormValue("sort")
	order := r.FormValue("order")

	// Fetch products from the in-memory list (replace with actual database query)
	products := searchInMemoryProducts(searchTerm, filters)

	// Sort products based on the specified criteria
	if sort == "price" {
		products = sortProductsByPrice(products, order)
	}

	// Respond with the search results
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// searchInMemoryProducts performs a search on an in-memory list of products
func searchInMemoryProducts(searchTerm string, filters map[string]string) []Product {
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
		for _, value := range filters {
			if strings.ToLower(product.Category) != strings.ToLower(value) {
				match = false
				break
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

func main() {
	http.HandleFunc("/search", SearchProducts)
	fmt.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
