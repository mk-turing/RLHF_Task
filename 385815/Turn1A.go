package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// TenantClaims represents the JWT claims for a tenant
type TenantClaims struct {
	jwt.StandardClaims
	TenantID string `json:"tenant_id"`
}

func main() {
	fmt.Println("Application started...")

	// Load the common secret key
	//secKeyPath := os.Getenv("JWT_SECRET_KEY_PATH")
	//if secKeyPath == "" {
	//	fmt.Println("JWT_SECRET_KEY_PATH environment variable not set.")
	//	os.Exit(1)
	//}
	//secKey, err := os.ReadFile(secKeyPath)
	//if err != nil {
	//	fmt.Println("Error reading secret key file:", err)
	//	os.Exit(1)
	//}
	secKey := []byte("your-secret-key")

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
		// For simplicity, we'll allow logins from any username in this example
		claims := &TenantClaims{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(1 * time.Hour).UnixMilli(),
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

		c.JSON(http.StatusOK, gin.H{
			"token": signedToken,
			"user":  user,
		})
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
