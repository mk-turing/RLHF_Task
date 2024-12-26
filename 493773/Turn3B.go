package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type User struct {
	ID       int     `json:"id"`
	Username string  `json:"username"`
	Posts    []*Post `json:"posts,omitempty"`
}

type Post struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	UserID int    `json:"user_id"`
}

type PaginatedResponse struct {
	Data  interface{} `json:"data"` // Can hold either []User or []Post
	Page  int         `json:"page"`
	Size  int         `json:"size"`
	Total int         `json:"total"`
	Sort  string      `json:"sort"`
}

// In-memory data for demonstration purposes
var users []*User
var posts []*Post

func init() {
	// Initialize user and post data
	for i := 1; i <= 100; i++ {
		user := &User{ID: i, Username: fmt.Sprintf("user%d", i)}
		users = append(users, user)
		for j := 1; j <= 10; j++ {
			post := &Post{ID: 10*i + j, Title: fmt.Sprintf("Post %d by User %d", 10*i+j, i), UserID: i}
			posts = append(posts, post)
			user.Posts = append(user.Posts, post)
		}
	}
}

func paginateUsersHandler(w http.ResponseWriter, r *http.Request) {
	// ... (Same as before)
	paginatedUsers := make([]*User, 0, size)
	for _, user := range users[start:end] {
		// If the user_id query parameter is provided, filter users by it
		userIDParam := r.URL.Query().Get("user_id")
		if userIDParam != "" {
			userID, _ := strconv.Atoi(userIDParam)
			if user.ID == userID {
				paginatedUsers = append(paginatedUsers, user)
			}
		} else {
			paginatedUsers = append(paginatedUsers, user)
		}
	}

	total := len(paginatedUsers)
	paginatedResponse := PaginatedResponse{
		Data:  paginatedUsers,
		Page:  page,
		Size:  size,
		Total: total,
		Sort:  sortBy,
	}
	// ... (Same as before)
}

func paginateHandler(w http.ResponseWriter, r *http.Request) {
	// ... (Same as before)

	// Check if the request is for users or posts
	isPostRequest := strings.HasPrefix(r.URL.Path, "/api/posts")

	var paginatedData interface{}
	var total int

	if isPostRequest {
		// Handle post pagination
		paginatedPosts := make([]*Post, 0, size)
		for _, post := range posts[start:end] {
			// Apply filters based on query parameters
			userIDParam := r.URL.Query().Get("user_id")
			if userIDParam != "" {
				userID, _ := strconv.Atoi(userIDParam)
				if post.UserID == userID {
					paginatedPosts = append(paginatedPosts, post)
				}
			} else {
				paginatedPosts = append(paginatedPosts, post)
			}