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
}

var items []Item

func init() {
	for i := 1; i <= 1000; i++ {
		items = append(items, Item{ID: i, Name: fmt.Sprintf("Item %d", i)})
	}
}

func paginateHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))

	// Default values
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10 // Default to 10 items per page
	}

	// Calculate start and end indices
	start := (page - 1) * size
	end := start + size
	if end > len(items) {
		end = len(items)
	}

	// Total number of items
	total := len(items)

	// Prepare the paginated response
	paginatedResponse := PaginatedResponse{
		Items: items[start:end],
		Page:  page,
		Size:  size,
		Total: total,
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
