package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"net/http"
	"strings"
)

var (
	clientID     = "your-client-id"
	clientSecret = "your-client-secret"
	redirectURL  = "http://localhost:8080/callback"
	oauth2Config = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     github.Endpoint,
		Scopes:       []string{"user:email"},
	}
)

func generateState() (string, error) {
	b := make([]byte, 32) // Generate a random state using 256 bits
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func login(c *gin.Context) {
	state, err := generateState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate state"})
		return
	}
	// Store the state in **session** or **temporary storage**

	url := oauth2Config.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func callback(c *gin.Context) {
	// CSRF check: Retrieve stored state and validate
	code := c.Query("code")
	token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not exchange token"})
		return
	}

	// Validate token and fetch user info (as shown before)

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged in!", "token": token.AccessToken})
}

func ValidateAccessTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
			c.Abort()
			return
		}
		// Verify the token and extract claims based on your OAuth provider
		// Ensure that token is not expired
		// If invalid, respond with an error
		// Placeholder logic to check token
		if tokenString != "expected-token" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
	}
}

func main() {
	router := gin.Default()

	// Use access token validation middleware for secured endpoints
	router.Use(ValidateAccessTokenMiddleware())

	router.GET("/login", login)
	router.GET("/callback", callback)

	// Refresh token endpoint
	router.POST("/refresh", func(c *gin.Context) {
		refreshToken := c.PostForm("refresh_token")

		// Validate the refresh token, check if it's stored in a database or memory
		// If valid, obtain a new access token

		tokenURL := "https://github.com/login/oauth/access_token"
		req, err := http.NewRequest("POST", tokenURL, strings.NewReader("grant_type=refresh_token&refresh_token="+refreshToken))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create request"})
			return
		}
		req.SetBasicAuth(clientID, clientSecret)

		// Perform the request to get new access token
		// Handle errors and validate response...
	})

	router.Run(":8080")
}
