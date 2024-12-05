package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func main() {
	// Define the signing key (in a real application, use a more secure method)
	signingKey := []byte("your_secret_signing_key")

	// Create a new JWT claim set
	claims := jwt.MapClaims{
		// Registered Claims
		"iss": "your_domain.com",
		"sub": "user123",
		"exp": time.Now().Add(time.Hour * 1).Unix(), // Token expires in 1 hour
		"nbf": time.Now().Unix(),                    // Token not before current time
		"iat": time.Now().Unix(),                    // Token issued at current time
		"aud": "api.your_domain.com",

		// Public Claim
		"user_id": "456",
		"name":    "John Doe",
	}

	// Create a new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token
	signedToken, err := token.SignedString(signingKey)
	if err != nil {
		panic(err)
	}

	fmt.Println("Generated Token:", signedToken)

	// Parse the token
	parsedToken, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the signature here in a real application
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return signingKey, nil
	})
	if err != nil {
		panic(err)
	}

	// Get claims from the parsed token
	claimsMap, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		panic("Invalid claim type")
	}

	// Access claims
	fmt.Println("Issuer:", claimsMap["iss"])
	fmt.Println("Subject:", claimsMap["sub"])
	fmt.Println("User ID:", claimsMap["user_id"])
	fmt.Println("Name:", claimsMap["name"])
}
