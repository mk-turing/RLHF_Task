package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/golang-jwt/jwt"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

var (
	mu         sync.RWMutex
	doc        *Document = &Document{Content: ""}
	connections map[*websocket.Conn]bool = make(map[*websocket.Conn]bool)
	users       map[int]*User = make(map[int]*User)
)

func main() {
	// Initialize sample users
	users[1] = &User{ID: 1, Username: "admin", Role: "admin"}
	users[2] = &User{ID: 2, Username: "editor", Role: "editor"}
	users[3] = &User{ID: 3, Username: "viewer", Role: "viewer"}

	r := mux.NewRouter()
	r.HandleFunc("/login", loginHandler).Methods("POST")
	r.HandleFunc("/ws", handleConnection).Methods("GET")
	r.HandleFunc("/document", getDocumentHandler).Methods("GET")

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

// Login handler for authentication
func loginHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate username and password (in memory in this example)
	user, ok := users[1] // Replace 1 with actual user ID based on username
	if !ok || user.Password != creds.Password {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}

// Handle WebSocket connections
func handleConnection(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}
	defer conn.Close()

	// Validate JWT token here