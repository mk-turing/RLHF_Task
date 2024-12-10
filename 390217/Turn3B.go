package main

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"

	"crypto/rand"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"` // Expiry in seconds
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	IssuedAt     int64  `json:"issued_at"`
}

var tokenCache = map[string]*Token{} // In production, use a more secure and persistent storage

func generateToken() (*Token, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	accessToken := base64.URLEncoding.EncodeToString(b)
	refreshToken := base64.URLEncoding.EncodeToString(b)

	token := &Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    3600, // 1 hour expiry for access token
		TokenType:    "Bearer",
		Scope:        "user:email",
		IssuedAt:     time.Now().Unix(),
	}

	return token, nil
}

func refreshToken(c *gin.Context) {
	refreshToken := c.PostForm("refresh_token")

	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing refresh_token"})
		return
	}

	token, found := tokenCache[refreshToken]

	if !found {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh_token"})
		return
	}

	// Check if access token has expired
	if time.Now().Unix()-token.IssuedAt < int64(token.ExpiresIn) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "access token has not expired"})
		return
	}

	// Generate a new access token
	newToken, err := generateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate new token"})
		return
	}

	delete(tokenCache, token.RefreshToken)
	tokenCache[newToken.RefreshToken] = newToken

	c.JSON(http.StatusOK, newToken)
}

var tokenBlacklist = map[string]bool{} // In production, use a persistent or distributed storage

func revokeToken(c *gin.Context) {
	token := c.PostForm("token")

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing token"})
		return
	}

	tokenBlacklist[token] = true

	c.JSON(http.StatusOK, gin.H{"message": "token revoked"})
}

func validateToken(token string) bool {
	_, found := tokenBlacklist[token]
	return !found
}

func main() {
	router := gin.Default()

	router.POST("/revoke", revokeToken)

	// Example middleware to validate tokens before accessing secured routes
	router.Use(func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token != "" {
			bearer := strings.Split(token, " ")
			if len(bearer) == 2 && bearer[0] == "Bearer" {
				accessToken := bearer[1]
				if !validateToken(accessToken) {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "token revoked"})
					c.Abort()
					return
				}
			}
		}
		c.Next()
	})

	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "protected resource"})
	})

	router.Run(":8080")
}
