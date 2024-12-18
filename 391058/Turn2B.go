package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"
)

// User represents a system user
type User struct {
	Username string
	Role     string
}

// InventoryUpdate represents an inventory update message
type InventoryUpdate struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

var (
	upGrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for simplicity
		},
	}
	// In-memory user store for demonstration purposes. Use a database in a real-world app
	users = map[string]User{
		"admin": {Username: "admin", Role: "admin"},
		"user":  {Username: "user", Role: "user"},
	}
	secretKey     = []byte("secret-key-goes-here") // Use a secure random key in a real app
	inventory    = make(map[string]int)
	inventoryMutex sync.Mutex
)

// GenerateJWT generates a JWT token for a given user
func GenerateJWT(user User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString(secretKey)
}

// VerifyJWT verifies a JWT token and extracts the user claims
func VerifyJWT(tokenString string) (User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return User{}, err
	}

	var user User
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user.Username = claims["username"].(string)
		user.Role = claims["role"].(string)
	} else {
		return User{}, fmt.Errorf("Invalid token claims")
	}

	return user, nil
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Validate the Authorization header to get the JWT token
	authHeader := r.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == "" {
		http.Error(w, "Invalid or missing Authorization header", http.StatusUnauthorized)
		return