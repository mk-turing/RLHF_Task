package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
)

type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type PaginatedResponse struct {
	Items []Item `json:"items"`
	Page  int    `json:"page"`
	Size  int    `json:"size"`
	Total int    `json:"total"`
	Sort  string `json:"sort"`
}

var items []Item

func init() {
	for i := 1; i <= 1000; i++ {
		items = append(items, Item{ID: i, Name: fmt.Sprintf("Item %d", i)})
	}
}

func paginateHandler(w http.ResponseWriter, r *http.Request) {
	// Log the request for debugging and monitoring
	log.Printf("%s %s\n", r.Method, r.RequestURI)

	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))
	sortBy := r.URL.Query().Get("sort")

	// Default values
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10 // Default to 10 items per page
	}

	// Validate sort parameter
	validSortParams := []string{"id", "name"}
	if sortBy != "" && !contains(validSortParams, sortBy) {
		http.Error(w, "Invalid sort parameter. Valid options are: id, name", http.StatusBadRequest)
		return
	}

	// Calculate start and end indices
	start := (page - 1) * size
	end := start + size
	if end > len(items) {
		end = len(items)
	}

	// Apply sorting if provided
	if sortBy != "" {
		sort.Slice(items[start:end], func(i, j int) bool {
			switch sortBy {
			case "id":
				return items[start+i].ID < items[start+j].ID
			case "name":
				return items[start+i].Name < items[start+j].Name
			}
			return false // This should never be reached
		})
	}

	// Total number of items
	total := len(items)

	// Prepare the paginated response
	paginatedResponse := PaginatedResponse{
		Items: items[start:end],
		Page:  page,
		Size:  size,
		Total: total,
		Sort:  sortBy,
	}

	// Encode and send the response
	json.NewEncoder(w).Encode(paginatedResponse)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func main() {
	http.HandleFunc("/api/items", paginateHandler)

	fmt.Println("Server is running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
