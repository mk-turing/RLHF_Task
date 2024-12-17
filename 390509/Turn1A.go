package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

// Post represents a blog post.
type Post struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Author    string `json:"author"`
	Timestamp string `json:"timestamp"`
}

// PostService interface for handling post operations.
type PostService interface {
	Create(Post) (*Post, error)
	GetAll() ([]Post, error)
	GetById(int) (*Post, error)
	Update(Post, int) (*Post, error)
	Delete(int) error
}

// PostManager struct implements the PostService interface.
type PostManager struct {
	posts []Post
}

// Create adds a new post.
func (pm *PostManager) Create(post Post) (*Post, error) {
	post.ID = len(pm.posts) + 1
	pm.posts = append(pm.posts, post)
	return &post, nil
}

// GetAll retrieves all posts.
func (pm *PostManager) GetAll() ([]Post, error) {
	return pm.posts, nil
}

// GetById retrieves a post by its ID.
func (pm *PostManager) GetById(id int) (*Post, error) {
	for _, post := range pm.posts {
		if post.ID == id {
			return &post, nil
		}
	}
	return nil, fmt.Errorf("Post not found with ID: %d", id)
}

// Update updates an existing post.
func (pm *PostManager) Update(post Post, id int) (*Post, error) {
	for i, p := range pm.posts {
		if p.ID == id {
			pm.posts[i] = post
			return &post, nil
		}
	}
	return nil, fmt.Errorf("Post not found with ID: %d", id)
}

// Delete removes a post by its ID.
func (pm *PostManager) Delete(id int) error {
	for i, p := range pm.posts {
		if p.ID == id {
			pm.posts = append(pm.posts[:i], pm.posts[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Post not found with ID: %d", id)
}

// PostController handles requests related to posts.
type PostController struct {
	service PostService
}

func NewPostController(service PostService) *PostController {
	return &PostController{service}
}

// CreatePost handler for creating a new post.
func (pc *PostController) CreatePost(w http.ResponseWriter, r *http.Request) {
	var post Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	newPost, err := pc.service.Create(post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newPost)
}

// GetAllPosts handler for retrieving all posts.
func (pc *PostController) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := pc.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}

// GetPostById handler for retrieving a post by its ID.
func (pc *PostController) GetPostById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	post, err := pc.service.GetById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(post)
}

// UpdatePost handler for updating an existing post.
func (pc *PostController) UpdatePost(w http.ResponseWriter, r *http.Request) {
	var post Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedPost, err := pc.service.Update(post, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedPost)
}

// DeletePost handler for deleting a post by its ID.
func (pc *PostController) DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := pc.service.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	postManager := &PostManager{}
	postController := NewPostController(postManager)

	r := mux.NewRouter()
	r.HandleFunc("/posts", postController.CreatePost).Methods("POST")
	r.HandleFunc("/posts", postController.GetAllPosts).Methods("GET")
	r.HandleFunc("/posts/{id}", postController.GetPostById).Methods("GET")
	r.HandleFunc("/posts/{id}", postController.UpdatePost).Methods("PUT")
	r.HandleFunc("/posts/{id}", postController.DeletePost).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), r))
}
