package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	Sort  string `json:"sort"` // New field to indicate sort order
}

var items []Item

func init() {
	for i := 1; i <= 1000; i++ {
		items = append(items, Item{ID: i, Name: fmt.Sprintf("Item %d", i)})
	}
}

func sortItems(items []Item, sortField string, isAscending bool) []Item {
	sortedItems := make([]Item, len(items))
	copy(sortedItems, items)

	switch sortField {
	case "id":
		if isAscending {
			return sortedItems
		} else {
			return sortedItems[::-1]
		}
	case "name":
		if isAscending {
			return sortedItems
		} else {
			return sortedItems[::-1]
		}
	default:
		return items
	}
}

func paginateHandler(w http.ResponseWriter, r *http.Request) {
	// Log the request for debugging and monitoring
	log.Printf("Request: %s", r.URL.RequestURI())

	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))
	sortField := r.URL.Query().Get("sort")
	sortOrder := r.URL.Query().Get("order")

	// Default values
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10 // Default to 10 items per page
	}

	isAscending := sortOrder == "asc" || sortOrder == ""
	if sortField != "" && sortField != "id" && sortField != "name" {
		http.Error(w, fmt.Sprintf("Invalid sort field: %s", sortField), http.StatusBadRequest)
		return
	}

	// Calculate start and end indices
	start := (page - 1) * size
	end := start + size
	if end > len(items) {
		end = len(items)
	}

	// Sort the items
	sortedItems := sortItems(items, sortField, isAscending)

	// Total number of items
	total := len(items)

	// Prepare the paginated response
	paginatedResponse := PaginatedResponse{
		Items: sortedItems[start:end],
		Page:  page,
		Size:  size,
		Total: total,
		Sort:  fmt.Sprintf("%s,%s", sortField, sortOrder),
	}

	// Encode and send the response
	json.NewEncoder(w).Encode(paginatedResponse)
}

func main() {
	http.HandleFunc("/api/items", paginateHandler)

	fmt.Println("Server is running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}