package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Role struct {
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

var roles = map[string]Role{
	"admin": {
		Name:        "admin",
		Permissions: []string{"read:all", "write:all", "delete:all"},
	},
	"editor": {
		Name:        "editor",
		Permissions: []string{"read:all", "write:all"},
	},
	"user": {
		Name:        "user",
		Permissions: []string{"read:all"},
	},
}

func getUserRole(userID string) string {
	// Mocking a function that retrieves a user's role by user ID
	// In practice, this could involve checking against a database or an external service
	if userID == "1" {
		return "admin"
	}
	return "user" // For simplicity, returning user for every other case
}

func authorize(allowedPermissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-User-ID") // Assuming we pass user's ID in headers
		userRole := getUserRole(userID)

		role, exists := roles[userRole]
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "invalid role"})
			c.Abort()
			return
		}

		for _, permission := range allowedPermissions {
			for _, rolePermission := range role.Permissions {
				if permission == rolePermission {
					// User has access
					c.Next()
					return
				}
			}
		}

		// User doesn't have required permission
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		c.Abort()
	}
}

func main() {

	router := gin.Default()

	// An example of a public route
	router.GET("/public", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "This is a public route."})
	})

	// Protected route that requires admin permissions
	router.GET("/admin", authorize("read:all", "write:all", "delete:all"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "This is an admin route."})
	})

	// Protected route that only needs user permissions
	router.GET("/user", authorize("read:all"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "This is a user route."})
	})

	router.Run(":8080")
}
