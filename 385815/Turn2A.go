package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

// TenantClaims represents the JWT claims for a tenant
type TenantClaims struct {
	jwt.StandardClaims
	TenantID string `json:"tenant_id"`
}

func main() {
	fmt.Println("Application started...")

	// Load the common secret key
	secKey := []byte("your-secret-key") // In production, ensure to load this securely

	router := gin.Default()

	// Sample endpoint to generate JWT for a given tenant
	router.POST("/login", func(c *gin.Context) {
		var user struct {
			Username string `json:"username"`
			TenantID string `json:"tenant_id"`
		}

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Generate JWT with tenant-specific claims
		claims := &TenantClaims{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(1 * time.Hour).Unix(), // Token expires in 1 hour
				Issuer:    "multi-tenant-go",
			},
			TenantID: user.TenantID,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signedToken, err := token.SignedString(secKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Generate a refresh token (UUID)
		refreshToken := uuid.New().String()

		c.JSON(http.StatusOK, gin.H{
			"token":         signedToken,
			"refresh_token": refreshToken,
			"user":          user,
		})
	})

	// Endpoint to refresh JWT using refresh token
	router.POST("/refresh", func(c *gin.Context) {
		var refreshTokenRequest struct {
			RefreshToken string `json:"refresh_token"`
		}

		if err := c.ShouldBindJSON(&refreshTokenRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate the refresh token (for simplicity, assume it's valid here)
		// In real applications, verify the refresh token from the database

		// Generate new access token
		newClaims := &TenantClaims{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
				Issuer:    "multi-tenant-go",
			},
			TenantID: "tenant123", // Replace with actual tenant ID
		}
		newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
		signedNewToken, err := newToken.SignedString(secKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": signedNewToken,
		})
	})

	// Middleware for Tenant Context
	router.Use(authenticateTenant)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}

// Middleware to authenticate and extract tenant information from JWT
func authenticateTenant(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization token"})
		c.Abort()
		return
	}

	parts := strings.SplitN(tokenString, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token format"})
		c.Abort()
		return
	}

	token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method)
		}
		return []byte("your-secret-key"), nil // Verify the JWT using the secret key
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	// Store the TenantID from the JWT claims in the context
	if claims, ok := token.Claims.(*TenantClaims); ok && token.Valid {
		c.Set("tenant_id", claims.TenantID)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	// Proceed to the next handler
	c.Next()
}
