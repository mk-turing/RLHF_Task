package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"golang.org/x/crypto/bcrypt"
)

// Your encrypted key (base64 encoded)
var key = []byte("your_32_byte_aes_key_1234567890abcdef")

// User struct to represent users with roles
type User struct {
	Username string
	Password []byte
	Role     string
}

// Sample users with their roles
var users = map[string]User{
	"admin": {
		Username: "admin",
		Password: bcrypt.MustGenerateFromPassword([]byte("password"), bcrypt.DefaultCost),
		Role:     "admin",
	},
	"user": {
		Username: "user",
		Password: bcrypt.MustGenerateFromPassword([]byte("password"), bcrypt.DefaultCost),
		Role:     "user",
	},
}

// Role to permissions map
var roleToPermissions = map[string][]string{
	"admin": {"read", "write"},
	"user":  {"read"},
}

// Middleware to authenticate and authorize users
func authenticateAndAuthorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Parse the basic auth credentials
		authParts := strings.SplitN(authHeader, " ", 2)
		if len(authParts) != 2 || strings.ToLower(authParts[0]) != "basic" {
			http.Error(w, "Invalid Authorization header", http.StatusBadRequest)
			return
		}

		credentials, err := base64.StdEncoding.DecodeString(authParts[1])
		if err != nil {
			http.Error(w, "Invalid Authorization header", http.StatusBadRequest)
			return
		}

		// Parse the username and password
		usernamePassword := strings.SplitN(string(credentials), ":", 2)
		if len(usernamePassword) != 2 {
			http.Error(w, "Invalid Authorization header", http.StatusBadRequest)
			return
		}

		username, password := usernamePassword[0], usernamePassword[1]

		// Authenticate the user
		user, ok := users[username]
		if !ok || err := bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Authorize the user for the requested route
		route := chi.URLParam(r, "route")
		requiredPermissions, ok := roleToPermissions[user.Role]
		if !ok {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		switch route {
		case "notifications":
			if !contains(requiredPermissions, "read") {
				http.Error(w, "Forbidden", http.StatusForbidden)