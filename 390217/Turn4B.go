package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"` // Assume role is returned by the OAuth provider
}

var users map[string]*User = make(map[string]*User) // In practice, store this in a database

// Mock function to retrieve user information from the OAuth provider
func getUserFromOAuth(accessToken string) (*User, error) {
	// Replace this with actual GitHub API call
	user := &User{
		Login: "john_doe",
		Name:  "John Doe",
		Email: "john@example.com",
		Role:  "user", // This would come from the provider
	}
	users[user.Login] = user
	return user, nil
}

// Mock OAuth flow
func login(c *gin.Context) {
	// This should be your actual OAuth login handling
	accessToken := "your_access_token"
	user, err := getUserFromOAuth(accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch user information"})
		return
	}
	// Store user in session or context (simplified for this example)
	c.Set("user", user)
	c.JSON(http.StatusOK, gin.H{"message": "logged in", "user": user})
}

func hasPermission(role string, permission string) bool {
	switch role {
	case "admin":
		return true // Admin has all permissions
	case "editor":
		return permission == "create" || permission == "view" || permission == "update"
	case "user":
		return permission == "view" || permission == "update"
	default:
		return false
	}
}

func requirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(*User)
		if !hasPermission(user.Role, permission) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func protectedEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "protected resource"})
}

func main() {
	router := gin.Default()

	router.GET("/login", login)

	// Endpoints with permissions
	router.GET("/protected", requirePermission("view"), protectedEndpoint)
	router.POST("/protected", requirePermission("create"), protectedEndpoint)
	router.PUT("/protected", requirePermission("update"), protectedEndpoint)
	router.DELETE("/protected", requirePermission("delete"), protectedEndpoint)

	router.Run(":8080")
}
