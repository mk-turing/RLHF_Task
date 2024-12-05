package main

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
	jwt.StandardClaims
	UserID  string `json:"user_id"`
	IsAdmin bool   `json:"is_admin"`
	// Private Claim
	PrivateData int `json:"private_data"`
}

func main() {
	// Define the secret key to sign the token
	secretKey := []byte("mySecretKey")

	// Create a new token with custom claims
	claims := &CustomClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "example.com",
			Subject:   "abak",
			Audience:  "client.app",
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserID:      "user123",
		IsAdmin:     true,
		PrivateData: 42,
	}

	// Create a Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign the token using the secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		fmt.Println("Error signing token:", err)
		return
	}
	fmt.Println("Generated Token:", tokenString)

	// Parse the token to retrieve claims
	tkn, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		fmt.Println("Error parsing token:", err)
		return
	}
	if claims, ok := tkn.Claims.(*CustomClaims); ok && tkn.Valid {
		fmt.Println("Parsed Claims:")
		fmt.Println("Issuer:", claims.Issuer)
		fmt.Println("Subject:", claims.Subject)
		fmt.Println("User ID:", claims.UserID)
		fmt.Println("Is Admin:", claims.IsAdmin)
		fmt.Println("Private Data:", claims.PrivateData)
	} else {
		fmt.Println("Invalid token claims")
	}
}
