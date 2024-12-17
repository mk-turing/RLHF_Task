package main

// Import required packages
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

// Post represents a blog post with versioning.
type Post struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Author    string `json:"author"`
	Timestamp string `json:"timestamp"`
	Version   int    `json:"version"`
}

// PostService interface for handling post operations with versioning.
type PostService interface {
	Create(Post) (*Post, error)
	GetAll() ([]Post, error)
	GetById(int) (*Post, error)
	Update(Post, int) (*Post, error)
	Delete(int) error
}

// PostManager struct implements the PostService interface.
type PostManager struct {
	posts map[int]Post
}

// Create adds a new post with version set to 1.
func (pm *PostManager) Create(post Post) (*Post, error) {
	post.ID = len(pm.posts) + 1
	post.Version = 1
	pm.posts[post.ID] = post
	return &post, nil
}

// GetAll retrieves all posts.
func (pm *PostManager) GetAll() ([]Post, error) {
	var posts []Post
	for _, post := range pm.posts {
		posts = append(posts, post)
	}
	return posts, nil
}

// GetById retrieves a post by its ID.
func (pm *PostManager) GetById(id int) (*Post, error) {
	if post, exists := pm.posts[id]; exists {
		return &post, nil
	}
	return nil, fmt.Errorf("Post not found with ID: %d", id)
}

// Update updates an existing post.
func (pm *PostManager) Update(post Post, id int) (*Post, error) {
	if _, exists := pm.posts[id]; exists {
		post.ID = id
		post.Version += 1
		pm.posts[id] = post
		return &post, nil
	}
	return nil, fmt.Errorf("Post not found with ID: %d", id)
}

// Delete removes a post by its ID.
func (pm *PostManager) Delete(id int) error {
	if _, exists := pm.posts[id]; exists {
		delete(pm.posts, id)
		return nil
	}
	return fmt.Errorf("Post not found with ID: %d", id)
}

// PostControllerV1 handles requests for API version 1.
type PostControllerV1 struct {
	service PostService
}

func NewPostControllerV1(service PostService) *PostControllerV1 {
	return &PostControllerV1{service}
}

func (pc *PostControllerV1) CreatePost(w http.ResponseWriter, r *http.Request) {
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

func (pc *PostControllerV1) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := pc.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}

func (pc *PostControllerV1) GetPostById(w http.ResponseWriter, r *http.Request) {
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

func (pc *PostControllerV1) UpdatePost(w http.ResponseWriter, r *http.Request) {
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

func (pc *PostControllerV1) DeletePost(w http.ResponseWriter, r *http.Request) {
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

// PostControllerV2 handles requests for API version 2.
type PostControllerV2 struct {
	service PostService
}

func NewPostControllerV2(service PostService) *PostControllerV2 {
	return &PostControllerV2{service}
}

func (pc *PostControllerV2) CreatePost(w http.ResponseWriter, r *http.Request) {
	var post Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	post.Version = 2 // Version 2 logic
	newPost, err := pc.service.Create(post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newPost)
}

func (pc *PostControllerV2) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := pc.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}

func (pc *PostControllerV2) GetPostById(w http.ResponseWriter, r *http.Request) {
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

func (pc *PostControllerV2) UpdatePost(w http.ResponseWriter, r *http.Request) {
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

func (pc *PostControllerV2) DeletePost(w http.ResponseWriter, r *http.Request) {
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
	postManager := &PostManager{posts: make(map[int]Post)}
	postService := PostService(postManager)

	r := mux.NewRouter()

	// Routes for API version 1
	v1 := r.PathPrefix("/v1").Subrouter()
	v1.HandleFunc("/posts", NewPostControllerV1(postService).CreatePost).Methods("POST")
	v1.HandleFunc("/posts", NewPostControllerV1(postService).GetAllPosts).Methods("GET")
	v1.HandleFunc("/posts/{id}", NewPostControllerV1(postService).GetPostById).Methods("GET")
	v1.HandleFunc("/posts/{id}", NewPostControllerV1(postService).UpdatePost).Methods("PUT")
	v1.HandleFunc("/posts/{id}", NewPostControllerV1(postService).DeletePost).Methods("DELETE")

	// Routes for API version 2
	v2 := r.PathPrefix("/v2").Subrouter()
	v2.HandleFunc("/posts", NewPostControllerV2(postService).CreatePost).Methods("POST")
	v2.HandleFunc("/posts", NewPostControllerV2(postService).GetAllPosts).Methods("GET")
	v2.HandleFunc("/posts/{id}", NewPostControllerV2(postService).GetPostById).Methods("GET")
	v2.HandleFunc("/posts/{id}", NewPostControllerV2(postService).UpdatePost).Methods("PUT")
	v2.HandleFunc("/posts/{id}", NewPostControllerV2(postService).DeletePost).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), r))
}
