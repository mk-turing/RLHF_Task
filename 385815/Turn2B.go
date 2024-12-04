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

// Constants for token expiry durations
const (
	AccessTokenExpiry  = time.Duration(15) * time.Minute
	RefreshTokenExpiry = time.Duration(24) * time.Hour
)

// TenantClaims represents the JWT claims for a tenant
type TenantClaims struct {
	jwt.StandardClaims
	TenantID string `json:"tenant_id"`
}

// Function to generate refresh token (UUID)
func generateRefreshToken() string {
	refreshToken := uuid.New().String()
	return refreshToken
}

// Function to generate access token for a tenant
func generateAccessTokenForTenant(tenantID string, secKey []byte) (string, error) {
	expiresAt := time.Now().Add(AccessTokenExpiry)
	claims := &TenantClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
			Issuer:    "multi-tenant-go",
		},
		TenantID: tenantID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secKey)
}

func main() {
	fmt.Println("Application started...")

	// Load the common secret key
	secKey := []byte("your-secret-key") // In production, load from environment

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

		// Generate Access Token
		signedToken, err := generateAccessTokenForTenant(user.TenantID, secKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Generate Refresh Token
		refreshToken := generateRefreshToken()

		// Return tokens
		c.JSON(http.StatusOK, gin.H{
			"accessToken":  signedToken,
			"refreshToken": refreshToken,
		})
	})

	// Endpoint to refresh JWT using refresh token
	router.POST("/refresh", func(c *gin.Context) {
		refreshToken := c.Query("refreshToken")

		// Validate refresh token (simplified)
		if refreshToken == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token required"})
			return
		}

		// Simulate tenant extraction from the refresh token (would be validated in real app)
		tenantID := "tenant123" // Extract tenantID from your refresh token in real cases

		// Generate new access token for tenant
		newAccessToken, err := generateAccessTokenForTenant(tenantID, secKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"accessToken": newAccessToken,
		})
	})

	// Middleware for Tenant Context
	router.Use(authenticateTenant)

	// Run on port 8080 (default)
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
		return []byte("your-secret-key"), nil
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

// Example of using the tenant info from context in a handler
func someHandler(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	// Use tenantID for tenant-specific data queries
	c.JSON(http.StatusOK, gin.H{"tenant_id": tenantID})
}
