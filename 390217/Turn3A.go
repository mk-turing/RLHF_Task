package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
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

	tokenStore = make(map[string]Token) // Simulated token store
)

type Token struct {
	AccessToken  string
	RefreshToken string
	Expiry       time.Time
}

func login(c *gin.Context) {
	code := c.Query("code")
	token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not exchange token"})
		return
	}

	userID := "someUniqueUserID" // Demo purpose; replace with actual user ID
	tokenStore[userID] = Token{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Expiry:       time.Now().Add(time.Hour), // Adjust as needed
	}

	c.JSON(http.StatusOK, gin.H{"access_token": token.AccessToken, "refresh_token": token.RefreshToken})
}

func refreshToken(c *gin.Context) {
	userID := "someUniqueUserID" // Obtain this from the session or token

	storedToken, exists := tokenStore[userID]
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	if time.Now().After(storedToken.Expiry) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token expired"})
		return
	}

	newToken, err := oauth2Config.TokenSource(context.Background(), &oauth2.Token{
		RefreshToken: storedToken.RefreshToken,
	}).Token()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not refresh token"})
		return
	}

	// Update token store
	tokenStore[userID] = Token{
		AccessToken:  newToken.AccessToken,
		RefreshToken: newToken.RefreshToken,
		Expiry:       time.Now().Add(time.Hour), // New expiry
	}

	c.JSON(http.StatusOK, gin.H{"access_token": newToken.AccessToken, "refresh_token": newToken.RefreshToken})
}

func logout(c *gin.Context) {
	userID := "someUniqueUserID" // Obtain this from the session or token
	delete(tokenStore, userID)
	c.JSON(http.StatusOK, gin.H{"message": "successfully logged out"})
}

func main() {
	router := gin.Default()
	router.GET("/login", login)
	router.POST("/refresh", refreshToken)
	router.POST("/logout", logout)

	router.Run(":8080")
}
