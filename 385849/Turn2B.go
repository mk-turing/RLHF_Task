package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// CustomClaims represent our custom claims for the JWT.
type CustomClaims struct {
	jwt.StandardClaims
	Admin bool   `json:"admin"`
	Role  string `json:"role"`
}

// secretKey should be replaced with a secure, unique key in a production environment.
var secretKey = []byte("yourSecretKey")

func jwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token missing"})
			c.Abort()
			return
		}

		// Remove the "Bearer " prefix from the token
		token = strings.TrimPrefix(token, "Bearer ")

		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unsupported signing method")
			}
			return secretKey, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		claims, ok := parsedToken.Claims.(*CustomClaims)
		if !ok || !parsedToken.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// Add the claims to the request context for use in handlers
		c.Set("user", claims)
		c.Next()
	}
}

func protectedHandler(c *gin.Context) {
	user := c.MustGet("user").(*CustomClaims)

	fmt.Println("Protected route accessed by:", user.Subject)
	fmt.Println("Is admin:", user.Admin)
	fmt.Println("Role:", user.Role)

	c.String(http.StatusOK, "Welcome to the protected route!")
}

func main() {
	router := gin.Default()
	router.Use(jwtMiddleware())

	// Protected route that requires a valid JWT with custom claims
	router.GET("/protected", protectedHandler)

	router.Run(":8080")
}
