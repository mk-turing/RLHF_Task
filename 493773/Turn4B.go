package main

import (
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
	"strconv"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Post struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	UserID int    `json:"user_id"`
}

type PaginatedResponse struct {
	Items []interface{} `json:"items"`
	Page  int           `json:"page"`
	Size  int           `json:"size"`
	Total int           `json:"total"`
	Sort  string        `json:"sort"`
}

var rdb *redis.Client // Redis client for caching

func init() {
	// Initialize Redis client
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // Redis password (empty if not set)
		DB:       0,               // Redis database number
	})

	// Sample data initialization (for demonstration purposes)
	// In a real application, this would be replaced by database queries
	// ...
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