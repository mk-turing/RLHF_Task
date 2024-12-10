package main

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/time/rate"
	"net/http"
	"strings"
	"time"
)

var (
	oauth2Config = &oauth2.Config{
		ClientID:     "your-client-id",
		ClientSecret: "your-client-secret",
		RedirectURL:  "http://localhost:8080/callback",
		Endpoint:     oauth2.Endpoint{}, // Replace with actual OAuth provider endpoint
		Scopes:       []string{"user:email"},
	}
)

var limiter = rate.NewLimiter(rate.Every(time.Minute), 100) // 100 requests per minute

func rateLimitMiddleware(c *gin.Context) {
	if !limiter.Allow() {
		c.AbortWithStatus(http.StatusTooManyRequests)
		fmt.Println("Rate limit exceeded")
		return
	}
	c.Next()
}

func main() {
	router := gin.Default()

	// Serve over HTTPS
	port := ":8443"
	certFile := "../390217/cert.pem"
	keyFile := "../390217/key.pem"
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	server := &http.Server{
		Addr:      port,
		Handler:   router,
		TLSConfig: tlsConfig,
	}

	// Rate limiting middleware
	router.Use(rateLimitMiddleware) // Use the rate limiting middleware

	// OAuth routes
	router.GET("/login", handleLogin)
	router.GET("/callback", handleCallback)
	router.POST("/revoke", handleRevoke)

	// Start HTTPS server
	fmt.Printf("Starting server on %s\n", port)
	if err := server.ListenAndServeTLS(certFile, keyFile); err != nil {
		panic(err)
	}
}

// 1. Generate State Parameter
func generateState() string {
	const stateSize = 32 // Random bytes
	buf := make([]byte, stateSize)
	_, err := rand.Read(buf)
	if err != nil {
		panic(fmt.Errorf("failed to generate state: %w", err))
	}
	return base64.URLEncoding.EncodeToString(buf)
}

// 2. Handle Login and Redirect to OAuth Provider
func handleLogin(c *gin.Context) {
	state := generateState()
	c.SetCookie("oauthState", state, 0, "/", "", false, true) // Secure flag set
	url := oauth2Config.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// 3. Handle OAuth Callback
func handleCallback(c *gin.Context) {
	// Validate Redirect URI
	providedRedirectURL := c.DefaultQuery("redirect_uri", "")
	validRedirectURL := "http://localhost:8080/callback" // Replace with your actual URL
	if !strings.HasPrefix(providedRedirectURL, validRedirectURL) {
		c.AbortWithStatus(http.StatusBadRequest)
		fmt.Println("Invalid redirect URI")
		return
	}

	// Validate State
	expectedState, err := c.Cookie("oauthState")
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		fmt.Println("Missing state cookie")
		return
	}
	c.SetCookie("oauthState", "", -1, "/", "", false, true) // Delete state cookie

	if c.DefaultQuery("state", "") != expectedState {
		c.AbortWithStatus(http.StatusBadRequest)
		fmt.Println("State mismatch")
		return
	}

	// Proceed with token exchange
	code := c.DefaultQuery("code", "")
	token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not exchange token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged in!", "token": token.AccessToken})
}

// 4. Verify Access Token
func verifyToken(c *gin.Context) error {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.Error(fmt.Errorf("Missing authorization token")) // Use fmt.Errorf for error
		return fmt.Errorf("unauthorized")
	}

	if !isValidToken(token) {
		c.Error(fmt.Errorf("Invalid authorization token")) // Use fmt.Errorf for error
		return fmt.Errorf("unauthorized")
	}

	return nil
}

// Example of validating token (you can replace with real token validation logic)
func isValidToken(token string) bool {
	// Implement actual token verification here
	return true // For demo purposes, always returns true
}

// 5. Handle Token Revocation
func handleRevoke(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	revokeURL := "https://oauth2.googleapis.com/revoke" // Replace with the provider's revoke URL
	client := &http.Client{}
	req, err := http.NewRequestWithContext(context.Background(), "POST", revokeURL, strings.NewReader(fmt.Sprintf("token=%s", token)))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Token revoked"})
}
