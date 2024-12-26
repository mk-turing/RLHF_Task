package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Posts []Post `json:"posts"`
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

var users []User
var posts []Post

func init() {
	// Sample data initialization
	for i := 1; i <= 10; i++ {
		users = append(users, User{ID: i, Name: fmt.Sprintf("User %d", i)})
	}
	for i := 1; i <= 100; i++ {
		posts = append(posts, Post{ID: i, Title: fmt.Sprintf("Post %d", i), UserID: i%10 + 1})
	}
	for i := range users {
		users[i].Posts = filterPostsByUserID(users[i].ID)
	}
}

func filterPostsByUserID(userID int) []Post {
	var filteredPosts []Post
	for _, post := range posts {
		if post.UserID == userID {
			filteredPosts = append(filteredPosts, post)
		}
	}
	return filteredPosts
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
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
	if end > len(users) {
		end = len(users)
	}

	// Apply sorting if provided
	if sortBy != "" {
		sort.Slice(users[start:end], func(i, j int) bool {
			switch sortBy {
			case "id":
				return users[start+i].ID < users[start+j].ID
			case "name":
				return users[start+i].Name < users[start+j].Name
			}
			return false // This should never be reached
		})
	}

	// Total number of items
	total := len(users)

	// Prepare the paginated response
	paginatedResponse := PaginatedResponse{
		Items: convertToInterfaceSlice(users[start:end]),
		Page:  page,
		Size:  size,
		Total: total,
		Sort:  sortBy,
	}

	// Encode and send the response
	json.NewEncoder(w).Encode(paginatedResponse)
}

func postsHandler(w http.ResponseWriter, r *http.Request) {
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
	validSortParams := []string{"id", "title"}
	if sortBy != "" && !contains(validSortParams, sortBy) {
		http.Error(w, "Invalid sort parameter. Valid options are: id, title", http.StatusBadRequest)
		return
	}

	// Calculate start and end indices
	start := (page - 1) * size
	end := start + size
	if end > len(posts) {
		end = len(posts)
	}

	// Apply sorting if provided
	if sortBy != "" {
		sort.Slice(posts[start:end], func(i, j int) bool {
			switch sortBy {
			case "id":
				return posts[start+i].ID < posts[start+j].ID
			case "title":
				return posts[start+i].Title < posts[start+j].Title
			}
			return false // This should never be reached
		})
	}

	// Total number of items
	total := len(posts)

	// Prepare the paginated response
	paginatedResponse := PaginatedResponse{
		Items: convertToInterfaceSlice(posts[start:end]),
		Page:  page,
		Size:  size,
		Total: total,
		Sort:  sortBy,
	}

	// Encode and send the response
	json.NewEncoder(w).Encode(paginatedResponse)
}

func convertToInterfaceSlice(slice interface{}) []interface{} {
	var result []interface{}
	switch v := slice.(type) {
	case []User:
		for _, u := range v {
			result = append(result, u)
		}
	case []Post:
		for _, p := range v {
			result = append(result, p)
		}
	}
	return result
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
	http.HandleFunc("/api/users", usersHandler)
	http.HandleFunc("/api/posts", postsHandler)

	fmt.Println("Server is running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
